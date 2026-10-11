// Package csscolorparser provides function for parsing CSS color string as defined in the W3C's CSS color module level 4.
package csscolorparser

import (
	"errors"
	"math"
	"strings"
)

var black = Color{0, 0, 0, 1}

var (
	ErrInvalidHex      = errors.New("invalid hex")
	ErrInvalidRgb      = errors.New("invalid rgb")
	ErrInvalidHsl      = errors.New("invalid hsl")
	ErrInvalidHwb      = errors.New("invalid hwb")
	ErrInvalidHsv      = errors.New("invalid hsv")
	ErrInvalidLab      = errors.New("invalid lab")
	ErrInvalidLch      = errors.New("invalid lch")
	ErrInvalidOklab    = errors.New("invalid oklab")
	ErrInvalidOklch    = errors.New("invalid oklch")
	ErrInvalidColor    = errors.New("invalid color")
	ErrInvalidFunction = errors.New("invalid function")
	ErrInvalidUnknown  = errors.New("invalid unknown")
)

// Parse parses CSS color string and returns, if successful, a Color.
func Parse(s string) (Color, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	if s == "transparent" {
		return Color{0, 0, 0, 0}, nil
	}

	// Predefined name / keyword

	c, ok := NamedColors[s]
	if ok {
		return Color{float64(c[0]) / 255, float64(c[1]) / 255, float64(c[2]) / 255, 1}, nil
	}

	// Hexadecimal

	if strings.HasPrefix(s, "#") {
		c, ok := parseHex(s[1:])
		if ok {
			return c, nil
		}
		return black, ErrInvalidHex
	}

	// Hexadecimal without prefix

	c2, ok2 := parseHex(s)
	if ok2 {
		return c2, nil
	}

	// Color functions

	op := strings.Index(s, "(")
	hasSuffix := strings.HasSuffix(s, ")")

	if !(op != -1 && hasSuffix) {
		return black, ErrInvalidUnknown
	}

	fname := strings.TrimSpace(s[:op])

	pp := newParamParser(s[op+1 : len(s)-1])
	pp.space()

	var color_func string

	if fname == "color" {
		v, ok := pp.value()
		if ok {
			color_func = v
		} else {
			return black, ErrInvalidColor
		}
		pp.space()
	}

	v0, ok0 := pp.value()
	sp0 := pp.commaOrSpace()
	v1, ok1 := pp.value()
	sp1 := pp.commaOrSpace()
	v2, ok2 := pp.value()

	if !(ok0 && ok1 && ok2 && sp0 && sp1) {
		return black, ErrInvalidFunction
	}

	sp2 := pp.space()
	end := pp.isEnd()
	sp3 := pp.commaOrSlash()

	v3, ok3 := pp.value()

	var alpha float64

	if end {
		alpha = 1.0
	} else if (sp2 || sp3) && ok3 {
		if !pp.isEnd() {
			return black, ErrInvalidFunction
		}

		v, ok, _ := parsePercentOrFloat(v3)
		if !ok {
			return black, ErrInvalidFunction
		}
		alpha = clamp(v)
	} else {
		return black, ErrInvalidFunction
	}

	if fname == "rgb" || fname == "rgba" {
		r, okR, pctR := parsePercentOrFloat(v0)
		g, okG, pctG := parsePercentOrFloat(v1)
		b, okB, pctB := parsePercentOrFloat(v2)

		if okR && okG && okB {
			if !pctR {
				r = r / 255
			}
			if !pctG {
				g = g / 255
			}
			if !pctB {
				b = b / 255
			}
			return Color{
				clamp(r),
				clamp(g),
				clamp(b),
				alpha,
			}, nil
		}

		return black, ErrInvalidRgb
	}

	if fname == "hsl" || fname == "hsla" {
		h, okH := parseAngle(v0)
		s, okS, _ := parsePercentOrFloat(v1)
		l, okL, _ := parsePercentOrFloat(v2)

		if okH && okS && okL {
			return FromHsl(h, s, l, alpha), nil
		}

		return black, ErrInvalidHsl
	}

	if fname == "hwb" || fname == "hwba" {
		H, okH := parseAngle(v0)
		W, okW, _ := parsePercentOrFloat(v1)
		B, okB, _ := parsePercentOrFloat(v2)

		if okH && okW && okB {
			return FromHwb(H, W, B, alpha), nil
		}

		return black, ErrInvalidHwb
	}

	if fname == "hsv" || fname == "hsva" {
		h, okH := parseAngle(v0)
		s, okS, _ := parsePercentOrFloat(v1)
		v, okV, _ := parsePercentOrFloat(v2)

		if okH && okS && okV {
			return FromHsv(h, s, v, alpha), nil
		}

		return black, ErrInvalidHsv
	}

	if fname == "oklab" {
		l, okL, _ := parsePercentOrFloat(v0)
		a, okA, fmtA := parsePercentOrFloat(v1)
		b, okB, fmtB := parsePercentOrFloat(v2)

		if okL && okA && okB {
			if fmtA {
				a = remap(a, -1.0, 1.0, -0.4, 0.4)
			}
			if fmtB {
				b = remap(b, -1.0, 1.0, -0.4, 0.4)
			}
			return FromOklab(math.Max(l, 0), a, b, alpha), nil
		}

		return black, ErrInvalidOklab
	}

	if fname == "oklch" {
		l, okL, _ := parsePercentOrFloat(v0)
		c, okC, fmtC := parsePercentOrFloat(v1)
		h, okH := parseAngle(v2)

		if okL && okC && okH {
			if fmtC {
				c = c * 0.4
			}
			return FromOklch(math.Max(l, 0), math.Max(c, 0), h*math.Pi/180, alpha), nil
		}

		return black, ErrInvalidOklch
	}

	if fname == "lab" {
		l, okL, fmtL := parsePercentOrFloat(v0)
		a, okA, fmtA := parsePercentOrFloat(v1)
		b, okB, fmtB := parsePercentOrFloat(v2)

		if okL && okA && okB {
			if fmtL {
				l = l * 100
			}
			if fmtA {
				a = remap(a, -1, 1, -125, 125)
			}
			if fmtB {
				b = remap(b, -1, 1, -125, 125)
			}
			return FromLab(math.Max(l, 0), a, b, alpha), nil
		}

		return black, ErrInvalidLab
	}

	if fname == "lch" {
		l, okL, fmtL := parsePercentOrFloat(v0)
		c, okC, fmtC := parsePercentOrFloat(v1)
		h, okH := parseAngle(v2)

		if okL && okC && okH {
			if fmtL {
				l = l * 100
			}
			if fmtC {
				c = c * 150
			}
			return FromLch(math.Max(l, 0), math.Max(c, 0), h*math.Pi/180, alpha), nil
		}

		return black, ErrInvalidLch
	}

	if fname == "color" {
		if color_func == "srgb" {
			r, okR, _ := parsePercentOrFloat(v0)
			g, okG, _ := parsePercentOrFloat(v1)
			b, okB, _ := parsePercentOrFloat(v2)

			if okR && okG && okB {
				return Color{
					clamp(r),
					clamp(g),
					clamp(b),
					alpha,
				}, nil
			}
		} else if color_func == "srgb-linear" {
			r, okR, _ := parsePercentOrFloat(v0)
			g, okG, _ := parsePercentOrFloat(v1)
			b, okB, _ := parsePercentOrFloat(v2)

			if okR && okG && okB {
				return FromLinearRGB(clamp(r), clamp(g), clamp(b), alpha), nil
			}
		}

		return black, ErrInvalidColor
	}

	return black, ErrInvalidFunction
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color

func parseHex(s string) (c Color, ok bool) {
	c.A = 1
	ok = true

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		}
		ok = false
		return 0
	}

	n := len(s)
	if n == 6 || n == 8 {
		c.R = float64(hexToByte(s[0])<<4+hexToByte(s[1])) / 255
		c.G = float64(hexToByte(s[2])<<4+hexToByte(s[3])) / 255
		c.B = float64(hexToByte(s[4])<<4+hexToByte(s[5])) / 255
		if n == 8 {
			c.A = float64(hexToByte(s[6])<<4+hexToByte(s[7])) / 255
		}
	} else if n == 3 || n == 4 {
		c.R = float64(hexToByte(s[0])*17) / 255
		c.G = float64(hexToByte(s[1])*17) / 255
		c.B = float64(hexToByte(s[2])*17) / 255
		if n == 4 {
			c.A = float64(hexToByte(s[3])*17) / 255
		}
	} else {
		ok = false
	}
	return
}
