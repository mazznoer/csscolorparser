package csscolorparser

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

// Inspired by https://github.com/deanm/css-color-parser-js

var black = color.RGBA{0, 0, 0, 255}

// Parse parses CSS color string and returns, if successful, a color.Color.
func Parse(s string) (color.Color, error) {
	input := s
	s = strings.TrimSpace(strings.ToLower(s))

	if s == "transparent" {
		return color.RGBA{0, 0, 0, 0}, nil
	}

	if s == "rebeccapurple" {
		return color.RGBA{102, 51, 153, 255}, nil
	}

	// Predefined name / keyword
	c, ok := colornames.Map[s]
	if ok {
		return c, nil
	}

	// Hexadecimal
	if strings.HasPrefix(s, "#") {
		c, ok := parseHex(s[1:])
		if ok {
			return c, nil
		}
		return black, fmt.Errorf("Invalid hex color, %s", input)
	}

	op := strings.Index(s, "(")

	if (op != -1) && strings.HasSuffix(s, ")") {
		fname := strings.TrimSpace(s[:op])
		alpha := 1.0
		okA := true
		s = s[op+1 : len(s)-1]
		s = strings.ReplaceAll(s, ",", " ")
		s = strings.ReplaceAll(s, "/", " ")
		params := strings.Fields(s)

		if fname == "rgb" || fname == "rgba" {
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("%s() format needs 3 or 4 parameters, %s", fname, input)
			}
			r, okR := parsePercentOr255(params[0])
			g, okG := parsePercentOr255(params[1])
			b, okB := parsePercentOr255(params[2])
			if len(params) == 4 {
				alpha, okA = parsePercentOrFloat(params[3])
			}
			if !okR || !okG || !okB || !okA {
				return black, fmt.Errorf("Wrong %s() components, %s", fname, input)
			}
			return color.NRGBA{
				uint8(clamp0_1(r) * 255),
				uint8(clamp0_1(g) * 255),
				uint8(clamp0_1(b) * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil

		} else if fname == "hsl" || fname == "hsla" {
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("%s() format needs 3 or 4 parameters, %s", fname, input)
			}
			h, okH := parseAngle(params[0])
			s, okS := parsePercentOrFloat(params[1])
			l, okL := parsePercentOrFloat(params[2])
			if len(params) == 4 {
				alpha, okA = parsePercentOrFloat(params[3])
			}
			if !okH || !okS || !okL || !okA {
				return black, fmt.Errorf("Wrong %s() components, %s", fname, input)
			}
			r, g, b := hslToRgb(normalizeAngle(h), clamp0_1(s), clamp0_1(l))
			return color.NRGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil

		} else if fname == "hwb" || fname == "hwba" {
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("hwb() format needs 3 or 4 parameters, %s", input)
			}
			H, okH := parseAngle(params[0])
			W, okW := parsePercentOrFloat(params[1])
			B, okB := parsePercentOrFloat(params[2])
			if len(params) == 4 {
				alpha, okA = parsePercentOrFloat(params[3])
			}
			if !okH || !okW || !okB || !okA {
				return black, fmt.Errorf("Wrong hwb() components, %s", input)
			}
			r, g, b := hwbToRgb(normalizeAngle(H), clamp0_1(W), clamp0_1(B))
			return color.NRGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil

		} else if fname == "hsv" || fname == "hsva" {
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("hsv() format needs 3 or 4 parameters, %s", input)
			}
			h, okH := parseAngle(params[0])
			s, okS := parsePercentOrFloat(params[1])
			v, okV := parsePercentOrFloat(params[2])
			if len(params) == 4 {
				alpha, okA = parsePercentOrFloat(params[3])
			}
			if !okH || !okS || !okV || !okA {
				return black, fmt.Errorf("Wrong hsv() components, %s", input)
			}
			r, g, b := hsvToRgb(normalizeAngle(h), clamp0_1(s), clamp0_1(v))
			return color.NRGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil
		}
	}

	c2, ok2 := parseHex(s)
	if ok2 {
		return c2, nil
	}

	return black, fmt.Errorf("Invalid color format, %s", input)
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color

func parseHex(s string) (c color.NRGBA, ok bool) {
	c.A = 255
	ok = true

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		ok = false
		return 0
	}

	n := len(s)
	if n == 6 || n == 8 {
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
		if n == 8 {
			c.A = hexToByte(s[6])<<4 + hexToByte(s[7])
		}
	} else if n == 3 || n == 4 {
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
		if n == 4 {
			c.A = hexToByte(s[3]) * 17
		}
	} else {
		ok = false
	}
	return
}

func hueToRgb(n1, n2, h float64) float64 {
	h = math.Mod(h, 6)
	if h < 1 {
		return n1 + ((n2 - n1) * h)
	}
	if h < 3 {
		return n2
	}
	if h < 4 {
		return n1 + ((n2 - n1) * (4 - h))
	}
	return n1
}

// h = 0..360
// s, l = 0..1
// r, g, b = 0..1
func hslToRgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}
	var n2 float64
	if l < 0.5 {
		n2 = l * (1 + s)
	} else {
		n2 = l + s - (l * s)
	}
	n1 := 2*l - n2
	h /= 60
	r = clamp0_1(hueToRgb(n1, n2, h+2))
	g = clamp0_1(hueToRgb(n1, n2, h))
	b = clamp0_1(hueToRgb(n1, n2, h-2))
	return
}

func hwbToRgb(hue, white, black float64) (r, g, b float64) {
	if white+black >= 1 {
		gray := white / (white + black)
		return gray, gray, gray
	}
	r, g, b = hslToRgb(hue, 1, 0.5)
	r = r*(1-white-black) + white
	g = g*(1-white-black) + white
	b = b*(1-white-black) + white
	return
}

func hsvToHsl(H, S, V float64) (h, s, l float64) {
	h = H
	s = S
	l = (2 - S) * V / 2
	if l != 0 {
		if l == 1 {
			s = 0
		} else if l < 0.5 {
			s = S * V / (l * 2)
		} else {
			s = S * V / (2 - l*2)
		}
	}
	return
}

func hsvToRgb(H, S, V float64) (r, g, b float64) {
	h, s, l := hsvToHsl(H, S, V)
	return hslToRgb(h, s, l)
}

func clamp0_1(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func parseFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f, err == nil
}

func parsePercentOrFloat(s string) (float64, bool) {
	if strings.HasSuffix(s, "%") {
		f, ok := parseFloat(s[:len(s)-1])
		if !ok {
			return 0, false
		}
		return f / 100, true
	}
	return parseFloat(s)
}

func parsePercentOr255(s string) (float64, bool) {
	if strings.HasSuffix(s, "%") {
		f, ok := parseFloat(s[:len(s)-1])
		if !ok {
			return 0, false
		}
		return f / 100, true
	}
	f, ok := parseFloat(s)
	if !ok {
		return 0, false
	}
	return f / 255, true
}

// Result angle in degrees (not normalized)
func parseAngle(s string) (float64, bool) {
	if strings.HasSuffix(s, "deg") {
		s = s[:len(s)-3]
		return parseFloat(s)
	}
	if strings.HasSuffix(s, "grad") {
		s = s[:len(s)-4]
		f, ok := parseFloat(s)
		if !ok {
			return 0, false
		}
		return f / 400 * 360, true
	}
	if strings.HasSuffix(s, "rad") {
		s = s[:len(s)-3]
		f, ok := parseFloat(s)
		if !ok {
			return 0, false
		}
		return f / math.Pi * 180, true
	}
	if strings.HasSuffix(s, "turn") {
		s = s[:len(s)-4]
		f, ok := parseFloat(s)
		if !ok {
			return 0, false
		}
		return f * 360, true
	}
	return parseFloat(s)
}

func normalizeAngle(t float64) float64 {
	t = math.Mod(t, 360)
	if t < 0 {
		t += 360
	}
	return t
}
