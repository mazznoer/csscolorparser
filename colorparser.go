package csscolorparser

import (
	"fmt"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

// Main algorithm based on: https://github.com/deanm/css-color-parser-js

var (
	transparent = color.RGBA{}
	black       = color.RGBA{0, 0, 0, 255}

	reHex       = regexp.MustCompile(`^#[0-9a-f]{3,}$`)
	reAngleDeg  = regexp.MustCompile(`^.+deg$`)
	reAngleGrad = regexp.MustCompile(`^.+grad$`)
	reAngleRad  = regexp.MustCompile(`^.+rad$`)
	reAngleTurn = regexp.MustCompile(`^.+turn$`)
)

func Parse(s string) (color.Color, error) {
	input := s
	s = strings.TrimSpace(strings.ToLower(s))

	if s == "transparent" {
		return transparent, nil
	}

	// Predefined name / keyword
	c, ok := colornames.Map[s]
	if ok {
		return c, nil
	}

	// Hexadecimal
	c2, ok := parseHex(s)
	if ok {
		return c2, nil
	}

	op := strings.Index(s, "(")
	ep := strings.Index(s, ")")

	if (op != -1) && (ep+1 == len(s)) {
		fname := s[:op]
		alpha := 1.0
		s = s[op+1 : ep]
		var params []string
		if strings.ContainsAny(s, ",") {
			s = strings.ReplaceAll(s, "/", ",")
			params = strings.Split(s, ",")
			for i, x := range params {
				params[i] = strings.TrimSpace(x)
			}
		} else {
			s = strings.ReplaceAll(s, "/", " ")
			params = strings.Fields(s)
		}

		switch fname {
		case "rgba":
			if len(params) != 4 {
				return black, fmt.Errorf("Invalid rgba() format, %v", input)
			}
			alpha, ok = parsePercentOrFloat(params[3])
			if !ok {
				return black, fmt.Errorf("Invalid rgba() format, %v", input)
			}
			params = params[:3]
			fallthrough

		case "rgb":
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("Invalid rgb() format, %v", input)
			}
			r, ok := parsePercentOr255(params[0])
			if !ok {
				return black, fmt.Errorf("Invalid rgb() format, %v", input)
			}
			g, ok := parsePercentOr255(params[1])
			if !ok {
				return black, fmt.Errorf("Invalid rgb() format, %v", input)
			}
			b, ok := parsePercentOr255(params[2])
			if !ok {
				return black, fmt.Errorf("Invalid rgb() format, %v", input)
			}
			if len(params) == 4 {
				alpha, ok = parsePercentOrFloat(params[3])
				if !ok {
					return black, fmt.Errorf("Invalid rgb() format, %v", input)
				}
			}
			return color.RGBA{
				uint8(clamp0_1(r) * 255),
				uint8(clamp0_1(g) * 255),
				uint8(clamp0_1(b) * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil

		case "hsla":
			if len(params) != 4 {
				return black, fmt.Errorf("Invalid hsla() format, %v", input)
			}
			alpha, ok = parsePercentOrFloat(params[3])
			if !ok {
				return black, fmt.Errorf("Invalid hsla() format, %v", input)
			}
			params = params[:3]
			fallthrough

		case "hsl":
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("Invalid hsl() format, %v", input)
			}
			h, ok := parseHue(params[0])
			if !ok {
				return black, fmt.Errorf("Invalid hsl() format, %v", input)
			}
			s, ok := parsePercentOrFloat(params[1])
			if !ok {
				return black, fmt.Errorf("Invalid hsl() format, %v", input)
			}
			l, ok := parsePercentOrFloat(params[2])
			if !ok {
				return black, fmt.Errorf("Invalid hsl() format, %v", input)
			}
			if len(params) == 4 {
				alpha, ok = parsePercentOrFloat(params[3])
				if !ok {
					return black, fmt.Errorf("Invalid hsl() format, %v", input)
				}
			}
			r, g, b := hslToRgb(normalizeAngle(h), clamp0_1(s), clamp0_1(l))
			return color.RGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				uint8(alpha * 255),
			}, nil

		case "hwb":
			if len(params) != 3 && len(params) != 4 {
				return black, fmt.Errorf("Invalid hwb() format, %v", input)
			}
			h, ok := parseHue(params[0])
			if !ok {
				return black, fmt.Errorf("Invalid hwb() format, %v", input)
			}
			w, ok := parsePercentOrFloat(params[1])
			if !ok {
				return black, fmt.Errorf("Invalid hwb() format, %v", input)
			}
			bk, ok := parsePercentOrFloat(params[2])
			if !ok {
				return black, fmt.Errorf("Invalid hwb() format, %v", input)
			}
			if len(params) == 4 {
				alpha, ok = parsePercentOrFloat(params[3])
				if !ok {
					return black, fmt.Errorf("Invalid hwb() format, %v", input)
				}
			}
			r, g, b := hwbToRgb(normalizeAngle(h), clamp0_1(w), clamp0_1(bk))
			return color.RGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				uint8(clamp0_1(alpha) * 255),
			}, nil

		default:
			return black, fmt.Errorf("Invalid color format, %v", input)
		}
	}

	return black, fmt.Errorf("Invalid color format, %v", input)
}

// parseHex taken from https://github.com/fogleman/colormap with some modification

func parseHex(x string) (color.Color, bool) {
	if !reHex.MatchString(x) {
		return transparent, false
	}
	var r, g, b, a int
	a = 255
	switch len(x) {
	case 4: // #rgb
		fmt.Sscanf(x, "#%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 5: // #rgba
		fmt.Sscanf(x, "#%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 7: // #rrggbb
		fmt.Sscanf(x, "#%02x%02x%02x", &r, &g, &b)
	case 9: // # rrggbbaa
		fmt.Sscanf(x, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	default:
		return transparent, false
	}
	return color.NRGBA64{
		uint16(r | r<<8),
		uint16(g | g<<8),
		uint16(b | b<<8),
		uint16(a | a<<8),
	}, true
}

func hue2rgb(n1, n2, h float64) float64 {
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

// h 0..360
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
	r = clamp0_1(hue2rgb(n1, n2, h+2))
	g = clamp0_1(hue2rgb(n1, n2, h))
	b = clamp0_1(hue2rgb(n1, n2, h-2))
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
	x := []rune(s)
	if string(x[len(x)-1]) == "%" {
		f, ok := parseFloat(s[:len(s)-1])
		if !ok {
			return 0, false
		}
		return f / 100, true
	}
	return parseFloat(s)
}

func parsePercentOr255(s string) (float64, bool) {
	x := []rune(s)
	if string(x[len(x)-1]) == "%" {
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
func parseHue(s string) (float64, bool) {
	if reAngleDeg.MatchString(s) {
		s = s[:len(s)-3]
		return parseFloat(s)
	}
	if reAngleGrad.MatchString(s) {
		s = s[:len(s)-4]
		f, ok := parseFloat(s)
		if !ok {
			return 0, false
		}
		return f / 400 * 360, true
	}
	if reAngleRad.MatchString(s) {
		s = s[:len(s)-3]
		f, ok := parseFloat(s)
		if !ok {
			return 0, false
		}
		return f / math.Pi * 180, true
	}
	if reAngleTurn.MatchString(s) {
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
