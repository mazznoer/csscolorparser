package csscolorparser

import "strings"

// Cursor-based scanner over an input string.
//
// All delimiters it cares about are ASCII, so it works on bytes and slicing
// is always safe for UTF-8 input (a delimiter byte can never be part of a
// multi-byte sequence).
type paramParser struct {
	s   string
	pos int
}

func newParamParser(s string) *paramParser {
	return &paramParser{s: s}
}

// value consumes and returns the substring from the current position up to
// (not including) the next space, comma, or slash. Those characters are
// ignored while inside parentheses, which may be nested.
//
// Returns ("", false) if there is no value at the current position; the
// position is left unchanged in that case.
func (p *paramParser) value() (string, bool) {
	start := p.pos
	i := start
	depth := 0

	for i < len(p.s) {
		c := p.s[i]
		if c == '(' {
			depth++
		} else if c == ')' {
			if depth > 0 {
				depth--
			}
		} else if depth == 0 && (isSpace(c) || c == ',' || c == '/') {
			break
		}
		i++
	}

	if i == start {
		return "", false
	}
	p.pos = i
	return p.s[start:i], true
}

// space consumes one or more ASCII whitespace characters.
// Returns true if at least one was found, false otherwise.
func (p *paramParser) space() bool {
	start := p.pos
	for p.pos < len(p.s) && isSpace(p.s[p.pos]) {
		p.pos++
	}
	return p.pos > start
}

// separator consumes a single character from chars, with optional
// whitespace around it. If no such character is found, nothing is consumed.
func (p *paramParser) separator(chars string) bool {
	start := p.pos
	p.space()
	if p.pos < len(p.s) && strings.IndexByte(chars, p.s[p.pos]) >= 0 {
		p.pos++
		p.space()
		return true
	}
	p.pos = start
	return false
}

// commaOrSpace consumes one or more spaces, or a single comma. Spaces are
// allowed around the comma. Returns true if either was found.
func (p *paramParser) commaOrSpace() bool {
	hadSpace := p.space()
	if p.pos < len(p.s) && p.s[p.pos] == ',' {
		p.pos++
		p.space()
		return true
	}
	return hadSpace
}

// commaOrSlash consumes a single comma or a single slash. Spaces are allowed
// around it. Returns true if one was found; otherwise nothing is consumed.
func (p *paramParser) commaOrSlash() bool {
	return p.separator(",/")
}

// slash consumes a single slash. Spaces are allowed around it. Returns true
// if a slash was found; otherwise nothing is consumed.
func (p *paramParser) slash() bool {
	return p.separator("/")
}

// isEnd reports whether the whole string has been read.
func (p *paramParser) isEnd() bool {
	return p.pos >= len(p.s)
}
