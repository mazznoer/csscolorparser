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

const maxRecDepth = 32

// Parse parses CSS color string and returns, if successful, a Color.
func Parse(s string) (Color, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	return parseAll(s, 0)
}

func parseAll(s string, depth int) (Color, error) {
	if depth > maxRecDepth {
		return black, ErrInvalidUnknown
	}

	// Parse absolute format

	c, err := parseAbs(s)
	if err == nil {
		return c, nil
	}

	// Parse relative format

	switch err {
	case ErrInvalidHex, ErrInvalidFunction, ErrInvalidUnknown:
		return black, err
	}

	op := strings.Index(s, "(")
	hasSuffix := strings.HasSuffix(s, ")")

	if !(op != -1 && hasSuffix) {
		return black, err
	}

	pp := newParamParser(s[op+1 : len(s)-1])
	pp.space()

	from, ok0 := pp.value()
	sp0 := pp.space()
	scol, ok1 := pp.value()
	sp1 := pp.space()

	if !(ok0 && ok1 && sp0 && sp1) {
		return black, err
	}

	if from != "from" {
		return black, err
	}

	col, errCol := parseAll(scol, depth+1)
	if errCol != nil {
		return black, err
	}

	var color_func string

	if err == ErrInvalidColor {
		v, ok := pp.value()
		if ok {
			color_func = v
		} else {
			return black, err
		}
		pp.space()
	}

	v0, okA := pp.value()
	spA := pp.space()
	v1, okB := pp.value()
	spB := pp.space()
	v2, okC := pp.value()

	if !(okA && okB && okC && spA && spB) {
		return black, err
	}

	pp.space()

	end := pp.isEnd()
	slash := pp.slash()

	v3, okD := pp.value()
	pp.space()

	var alpha string

	if end {
		alpha = "alpha"
	} else if slash && okD && pp.isEnd() {
		alpha = v3
	} else {
		return black, err
	}

	if err == ErrInvalidRgb {
		variables := [4]variable{
			{"r", col.R * 255},
			{"g", col.G * 255},
			{"b", col.B * 255},
			{"alpha", col.A},
		}
		r, okR, pctR := parseComponent(v0, variables)
		g, okG, pctG := parseComponent(v1, variables)
		b, okB, pctB := parseComponent(v2, variables)
		a, okA, _ := parseComponent(alpha, variables)

		if okR && okG && okB && okA {
			if !pctR {
				r = r / 255
			}
			if !pctG {
				g = g / 255
			}
			if !pctB {
				b = b / 255
			}
			return Color{r, g, b, a}, nil
		}

	} else if err == ErrInvalidHsl {
		hsl := col.Hsl()
		variables := [4]variable{
			{"h", hsl[0]},
			{"s", hsl[1] * 100},
			{"l", hsl[2] * 100},
			{"alpha", col.A},
		}
		h, okH, pctH := parseComponent(v0, variables)
		s, okS, pctS := parseComponent(v1, variables)
		l, okL, pctL := parseComponent(v2, variables)
		a, okA, _ := parseComponent(alpha, variables)

		if okH && okS && okL && okA && !pctH {
			if !pctS {
				s = s / 100
			}
			if !pctL {
				l = l / 100
			}
			return FromHsl(h, clamp(s), clamp(l), a), nil
		}

	} else if err == ErrInvalidHwb {
		hwb := col.Hwb()
		variables := [4]variable{
			{"h", hwb[0]},
			{"w", hwb[1] * 100},
			{"b", hwb[2] * 100},
			{"alpha", col.A},
		}
		h, okH, pctH := parseComponent(v0, variables)
		w, okW, pctW := parseComponent(v1, variables)
		b, okB, pctB := parseComponent(v2, variables)
		a, okA, _ := parseComponent(alpha, variables)

		if okH && okW && okB && okA && !pctH {
			if !pctW {
				w = w / 100
			}
			if !pctB {
				b = b / 100
			}
			return FromHwb(h, w, b, a), nil
		}

	} else if err == ErrInvalidLab {
		lab := col.Lab()
		variables := [4]variable{
			{"l", lab[0]},
			{"a", lab[1]},
			{"b", lab[2]},
			{"alpha", col.A},
		}
		l, okL, pctL := parseComponent(v0, variables)
		a, okA, pctA := parseComponent(v1, variables)
		b, okB, pctB := parseComponent(v2, variables)
		ap, okAp, _ := parseComponent(alpha, variables)

		if okL && okA && okB && okAp {
			if pctL {
				l *= 100
			}
			if pctA {
				a = remap(a, -1, 1, -125, 125)
			}
			if pctB {
				b = remap(b, -1, 1, -125, 125)
			}
			return FromLab(math.Max(l, 0), a, b, ap), nil
		}

	} else if err == ErrInvalidLch {
		lch := col.Lch()
		variables := [4]variable{
			{"l", lch[0]},
			{"c", lch[1]},
			{"h", lch[2] * (180 / math.Pi)},
			{"alpha", col.A},
		}
		l, okL, pctL := parseComponent(v0, variables)
		c, okC, pctC := parseComponent(v1, variables)
		h, okH, pctH := parseComponent(v2, variables)
		a, okA, _ := parseComponent(alpha, variables)

		if okL && okC && okH && okA && !pctH {
			if pctL {
				l *= 100
			}
			if pctC {
				c *= 150
			}
			return FromLch(math.Max(l, 0), math.Max(c, 0), h*(math.Pi/180), a), nil
		}

	} else if err == ErrInvalidOklab {
		lab := col.Oklab()
		variables := [4]variable{
			{"l", lab[0]},
			{"a", lab[1]},
			{"b", lab[2]},
			{"alpha", col.A},
		}
		l, okL, _ := parseComponent(v0, variables)
		a, okA, pctA := parseComponent(v1, variables)
		b, okB, pctB := parseComponent(v2, variables)
		ap, okAp, _ := parseComponent(alpha, variables)

		if okL && okA && okB && okAp {
			if pctA {
				a = remap(a, -1, 1, -0.4, 0.4)
			}
			if pctB {
				b = remap(b, -1, 1, -0.4, 0.4)
			}
			return FromOklab(math.Max(l, 0), a, b, ap), nil
		}

	} else if err == ErrInvalidOklch {
		lch := col.Oklch()
		variables := [4]variable{
			{"l", lch[0]},
			{"c", lch[1]},
			{"h", lch[2] * (180 / math.Pi)},
			{"alpha", col.A},
		}
		l, okL, _ := parseComponent(v0, variables)
		c, okC, pctC := parseComponent(v1, variables)
		h, okH, pctH := parseComponent(v2, variables)
		a, okA, _ := parseComponent(alpha, variables)

		if okL && okC && okH && okA && !pctH {
			if pctC {
				c *= 0.4
			}
			return FromOklch(math.Max(l, 0), math.Max(c, 0), h*(math.Pi/180), a), nil
		}

	} else if err == ErrInvalidColor {
		if color_func == "srgb" {
			variables := [4]variable{
				{"r", col.R},
				{"g", col.G},
				{"b", col.B},
				{"alpha", col.A},
			}
			r, okR, _ := parseComponent(v0, variables)
			g, okG, _ := parseComponent(v1, variables)
			b, okB, _ := parseComponent(v2, variables)
			a, okA, _ := parseComponent(alpha, variables)

			if okR && okG && okB && okA {
				return Color{r, g, b, a}, nil
			}

		} else if color_func == "srgb-linear" {
			rgb := col.LinearRgb()
			variables := [4]variable{
				{"r", rgb[0]},
				{"g", rgb[1]},
				{"b", rgb[2]},
				{"alpha", col.A},
			}
			r, okR, _ := parseComponent(v0, variables)
			g, okG, _ := parseComponent(v1, variables)
			b, okB, _ := parseComponent(v2, variables)
			a, okA, _ := parseComponent(alpha, variables)

			if okR && okG && okB && okA {
				return FromLinearRGB(r, g, b, a), nil
			}
		}
	}

	return black, err
}

func parseAbs(s string) (Color, error) {
	if s == "transparent" {
		return Color{0, 0, 0, 0}, nil
	}

	// Predefined name / keyword

	if c, ok := NamedColors[s]; ok {
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

	if c, ok := parseHex(s); ok {
		return c, nil
	}

	// Color functions

	op := strings.Index(s, "(")
	hasSuffix := strings.HasSuffix(s, ")")

	if !(op != -1 && hasSuffix) {
		return black, ErrInvalidUnknown
	}

	fname := strings.TrimSpace(s[:op])

	var err error
	if fname == "rgb" || fname == "rgba" {
		err = ErrInvalidRgb
	} else if fname == "hsl" || fname == "hsla" {
		err = ErrInvalidHsl
	} else if fname == "hwb" || fname == "hwba" {
		err = ErrInvalidHwb
	} else if fname == "hsv" {
		err = ErrInvalidHsv
	} else if fname == "lab" {
		err = ErrInvalidLab
	} else if fname == "lch" {
		err = ErrInvalidLch
	} else if fname == "oklab" {
		err = ErrInvalidOklab
	} else if fname == "oklch" {
		err = ErrInvalidOklch
	} else if fname == "color" {
		err = ErrInvalidColor
	} else {
		return black, ErrInvalidFunction
	}

	pp := newParamParser(s[op+1 : len(s)-1])
	pp.space()

	var color_func string

	if fname == "color" {
		v, ok := pp.value()
		if ok {
			color_func = v
		} else {
			return black, err
		}
		pp.space()
	}

	v0, ok0 := pp.value()
	sp0 := pp.commaOrSpace()
	v1, ok1 := pp.value()
	sp1 := pp.commaOrSpace()
	v2, ok2 := pp.value()

	if !(ok0 && ok1 && ok2 && sp0 && sp1) {
		return black, err
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
			return black, err
		}

		v, ok, _ := parsePercentOrFloat(v3)
		if !ok {
			return black, err
		}
		alpha = clamp(v)
	} else {
		return black, err
	}

	if err == ErrInvalidRgb {
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

	} else if err == ErrInvalidHsl {
		h, okH := parseAngle(v0)
		s, okS, pctS := parsePercentOrFloat(v1)
		l, okL, pctL := parsePercentOrFloat(v2)

		if okH && okS && okL {
			if !pctS {
				s = s / 100
			}
			if !pctL {
				l = l / 100
			}
			return FromHsl(h, s, l, alpha), nil
		}

	} else if err == ErrInvalidHwb {
		h, okH := parseAngle(v0)
		w, okW, pctW := parsePercentOrFloat(v1)
		b, okB, pctB := parsePercentOrFloat(v2)

		if okH && okW && okB {
			if !pctW {
				w = w / 100
			}
			if !pctB {
				b = b / 100
			}
			return FromHwb(h, w, b, alpha), nil
		}

	} else if err == ErrInvalidHsv {
		h, okH := parseAngle(v0)
		s, okS, _ := parsePercentOrFloat(v1)
		v, okV, _ := parsePercentOrFloat(v2)

		if okH && okS && okV {
			return FromHsv(h, s, v, alpha), nil
		}

	} else if err == ErrInvalidOklab {
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

	} else if err == ErrInvalidOklch {
		l, okL, _ := parsePercentOrFloat(v0)
		c, okC, fmtC := parsePercentOrFloat(v1)
		h, okH := parseAngle(v2)

		if okL && okC && okH {
			if fmtC {
				c = c * 0.4
			}
			return FromOklch(math.Max(l, 0), math.Max(c, 0), h*math.Pi/180, alpha), nil
		}

	} else if err == ErrInvalidLab {
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

	} else if err == ErrInvalidLch {
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

	} else if err == ErrInvalidColor {
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

		} else if color_func == "xyz-d50" {
			x, okX, _ := parsePercentOrFloat(v0)
			y, okY, _ := parsePercentOrFloat(v1)
			z, okZ, _ := parsePercentOrFloat(v2)

			if okX && okY && okZ {
				return FromXyzD50(x, y, z, alpha), nil
			}

		} else if color_func == "xyz-d65" || color_func == "xyz" {
			x, okX, _ := parsePercentOrFloat(v0)
			y, okY, _ := parsePercentOrFloat(v1)
			z, okZ, _ := parsePercentOrFloat(v2)

			if okX && okY && okZ {
				return FromXyzD65(x, y, z, alpha), nil
			}
		}
	}

	return black, err
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
