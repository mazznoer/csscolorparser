package csscolorparser

import (
	"fmt"
	"iter"
	"strings"
)

// ParseColors returns an iterator that lazily parses a comma-separated list
// of CSS colors.
//
// Commas inside parentheses are not treated as separators, so functional
// notations such as "rgb(255, 0, 0)" or "hsl(120, 100%, 50%)" stay intact.
// Empty entries (e.g. from a trailing comma) are skipped.
//
// Each iteration yields the parsed color and an error. If a color fails to
// parse, the error includes the offending token and iteration continues,
// letting the caller decide whether to stop (break) or skip.
//
// Example:
//
//	for c, err := range csscolorparser.ParseColors("red, #00ff00, rgb(0, 0, 255)") {
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		fmt.Println(c.HexString())
//	}
func ParseColors(s string) iter.Seq2[Color, error] {
	return func(yield func(Color, error) bool) {
		for token := range splitTopLevel(s) {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}

			c, err := Parse(token)
			if err != nil {
				err = fmt.Errorf("parse color %q: %w", token, err)
			}

			if !yield(c, err) {
				return
			}
		}
	}
}

// splitTopLevel yields the parts of s separated by commas that are not
// nested inside parentheses.
func splitTopLevel(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		depth, start := 0, 0

		for i := 0; i < len(s); i++ {
			switch s[i] {
			case '(':
				depth++
			case ')':
				if depth > 0 {
					depth--
				}
			case ',':
				if depth == 0 {
					if !yield(s[start:i]) {
						return
					}
					start = i + 1
				}
			}
		}

		yield(s[start:])
	}
}
