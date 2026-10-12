package csscolorparser

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type tokenKind uint8

const (
	tokEOF tokenKind = iota
	tokNum
	tokIdent
	tokAdd
	tokSub
	tokMul
	tokDiv
	tokLParen
	tokRParen
	tokUnknown
)

// token is a lexical token. Only the field matching kind is meaningful.
type token struct {
	kind  tokenKind
	num   float64 // tokNum
	ident string  // tokIdent
}

type lexer struct {
	src string
	pos int
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\f' || b == '\r'
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func isAlpha(b byte) bool { return (b >= 'a' && b <= 'z') }

// scan returns the next token. At end of input it returns tokEOF, and keeps
// returning it on every later call.
func (l *lexer) scan() token {
	for l.pos < len(l.src) && isSpace(l.src[l.pos]) {
		l.pos++
	}
	if l.pos >= len(l.src) {
		return token{kind: tokEOF}
	}

	b := l.src[l.pos]
	l.pos++

	switch b {
	case '+':
		return token{kind: tokAdd}
	case '-':
		return token{kind: tokSub}
	case '*':
		return token{kind: tokMul}
	case '/':
		return token{kind: tokDiv}
	case '(':
		return token{kind: tokLParen}
	case ')':
		return token{kind: tokRParen}
	}

	if isDigit(b) || b == '.' {
		return l.number(b)
	}
	if isAlpha(b) {
		return l.ident()
	}
	return token{kind: tokUnknown}
}

// number is called after the first byte (a digit or '.') has been consumed.
func (l *lexer) number(first byte) token {
	start := l.pos - 1
	dotSeen := first == '.'

	for l.pos < len(l.src) {
		c := l.src[l.pos]
		if isDigit(c) {
			l.pos++
		} else if c == '.' && !dotSeen {
			dotSeen = true
			l.pos++
		} else {
			break
		}
	}

	// A digit string too long for float64 comes back as ±Inf with ErrRange;
	// that's accepted here, and callers reject non-finite results.
	val, err := strconv.ParseFloat(l.src[start:l.pos], 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return token{kind: tokUnknown} // Triggered by isolated dots (e.g., ".")
	}
	return token{kind: tokNum, num: val}
}

// ident is called after the first byte (an ASCII letter) has been consumed.
// Identifiers consist solely of ASCII letters: no digits, no underscores.
func (l *lexer) ident() token {
	start := l.pos - 1
	for l.pos < len(l.src) && isAlpha(l.src[l.pos]) {
		l.pos++
	}
	return token{kind: tokIdent, ident: l.src[start:l.pos]}
}

//---

const maxDepth = 64

// prefixBP is the binding power of unary + and -, which binds tighter than
// any binary operator.
const prefixBP = 5

// VarLookup resolves a variable name. ok is false if the name is undefined.
type varLookup func(name string) (value float64, ok bool)

// Eval evaluates an arithmetic expression made of numbers, variables,
// parentheses, unary + and -, and the binary operators + - * /.
//
// Variable names consist only of ASCII letters and are resolved with vars;
// a nil vars means no variables are defined.
//
// ok is false on a syntax error, an undefined variable, nesting deeper than
// the internal limit, or trailing input after a complete expression.
// Division follows IEEE 754, so the result may be ±Inf or NaN.
func eval(src string, vars varLookup) (value float64, ok bool) {
	p := parser{lex: lexer{src: src}, vars: vars}
	p.next()

	v, ok := p.expr(0)
	if !ok || p.tok.kind != tokEOF {
		return 0, false
	}
	return v, true
}

type parser struct {
	lex   lexer
	tok   token // current lookahead; kind is tokEOF at end of input
	vars  varLookup
	depth int
}

func (p *parser) next() { p.tok = p.lex.scan() }

// infixBP returns the left and right binding powers of a binary operator.
// ok is false if kind is not a binary operator.
func infixBP(kind tokenKind) (left, right int, ok bool) {
	switch kind {
	case tokAdd, tokSub:
		return 1, 2, true
	case tokMul, tokDiv:
		return 3, 4, true
	}
	return 0, 0, false
}

// expr parses and evaluates an expression whose operators all bind at least
// as tightly as minBP.
func (p *parser) expr(minBP int) (float64, bool) {
	if p.depth > maxDepth {
		return 0, false
	}
	p.depth++
	defer func() { p.depth-- }()

	left, ok := p.operand()
	if !ok {
		return 0, false
	}

	for {
		lbp, rbp, isOp := infixBP(p.tok.kind)
		if !isOp || lbp < minBP {
			// EOF, a non-operator token, or an operator that binds too
			// loosely for this level. The caller decides what's valid.
			return left, true
		}

		op := p.tok.kind
		p.next()

		right, ok := p.expr(rbp)
		if !ok {
			return 0, false
		}

		switch op {
		case tokAdd:
			left += right
		case tokSub:
			left -= right
		case tokMul:
			left *= right
		case tokDiv:
			left /= right
		}
	}
}

// operand parses a number, variable, parenthesized expression, or a
// unary + / - applied to an operand.
func (p *parser) operand() (float64, bool) {
	switch p.tok.kind {
	case tokNum:
		v := p.tok.num
		p.next()
		return v, true

	case tokIdent:
		v, found := p.vars(p.tok.ident)
		if !found {
			return 0, false
		}
		p.next()
		return v, true

	case tokAdd, tokSub:
		negate := p.tok.kind == tokSub
		p.next()
		v, ok := p.expr(prefixBP)
		if !ok {
			return 0, false
		}
		if negate {
			v = -v
		}
		return v, true

	case tokLParen:
		p.next()
		v, ok := p.expr(0)
		if !ok || p.tok.kind != tokRParen {
			return 0, false
		}
		p.next()
		return v, true
	}

	return 0, false
}

//---

// Variable is a named constant usable as a bare keyword or inside calc(...).
type variable struct {
	Name  string
	Value float64
}

// parseComponent parses one component, which is one of:
//
//   - "none": zero, not a percentage
//   - a finite number, or a variable name
//   - a percentage like "50%", stored as 0.5 with IsPercent set
//   - calc(expr), where expr is evaluated by Eval and must come out finite
//
// s and the variable names are expected to be lowercase; matching is exact.

// returns (t, ok, pct)
func parseComponent(s string, variables [4]variable) (float64, bool, bool) {
	lookup := func(name string) (float64, bool) {
		if v, ok := parseNumber(name); ok {
			if !isFinite(v) {
				return 0, false
			}
			return v, true
		}
		for _, v := range variables {
			if name == v.Name {
				return v.Value, true
			}
		}
		return 0, false
	}

	// none
	if s == "none" {
		return 0, true, false
	}

	// bare number or keyword
	if v, ok := lookup(s); ok {
		return v, true, false
	}

	// percentage
	if rest, ok := strings.CutSuffix(s, "%"); ok {
		v, ok := parseNumber(rest)
		if !ok || !isFinite(v) {
			return 0, false, false
		}
		return v / 100, true, true
	}

	// calc(...)
	if rest, ok := strings.CutPrefix(s, "calc("); ok {
		if inner, ok := strings.CutSuffix(rest, ")"); ok {
			if v, ok := eval(inner, lookup); ok && isFinite(v) {
				return v, true, false
			}
		}
	}

	// invalid value
	return 0, false, false
}

// parseNumber parses a decimal floating-point number. strconv.ParseFloat is
// more permissive than wanted here (it accepts digit-separating underscores
// and hex floats such as "0x1p3"), so those forms are rejected up front.
// An out-of-range value parses successfully as ±Inf (strconv.ErrRange)
// instead of failing; callers reject non-finite values.
func parseNumber(s string) (float64, bool) {
	if strings.Contains(s, "_") {
		return 0, false
	}
	t := s
	if len(t) > 0 && (t[0] == '+' || t[0] == '-') {
		t = t[1:]
	}
	if len(t) >= 2 && t[0] == '0' && t[1] == 'x' { // input is lowercase
		return 0, false
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return 0, false
	}
	return v, true
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
