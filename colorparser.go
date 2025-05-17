// Package csscolorparser provides function for parsing CSS color string as defined in the W3C's CSS color module level 4.
package csscolorparser

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Inspired by https://github.com/deanm/css-color-parser-js

// R, G, B, A values in the range 0..1
type Color struct {
	R, G, B, A float64
}

// Implement the Go color.Color interface.
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R*c.A*65535 + 0.5)
	g = uint32(c.G*c.A*65535 + 0.5)
	b = uint32(c.B*c.A*65535 + 0.5)
	a = uint32(c.A*65535 + 0.5)
	return
}

// RGBA255 returns R, G, B, A values in the range 0..255
func (c Color) RGBA255() (r, g, b, a uint8) {
	r = uint8(c.R*255 + 0.5)
	g = uint8(c.G*255 + 0.5)
	b = uint8(c.B*255 + 0.5)
	a = uint8(c.A*255 + 0.5)
	return
}

// Clamp restricts R, G, B, A values to the range 0..1.
func (c Color) Clamp() Color {
	return Color{
		R: math.Max(math.Min(c.R, 1), 0),
		G: math.Max(math.Min(c.G, 1), 0),
		B: math.Max(math.Min(c.B, 1), 0),
		A: math.Max(math.Min(c.A, 1), 0),
	}
}

// HexString returns CSS hexadecimal string.
func (c Color) HexString() string {
	r, g, b, a := c.RGBA255()
	if a < 255 {
		return fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
	}
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// RGBString returns CSS RGB string.
func (c Color) RGBString() string {
	r, g, b, _ := c.RGBA255()
	if c.A < 1 {
		return fmt.Sprintf("rgba(%d,%d,%d,%v)", r, g, b, c.A)
	}
	return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
}

// Name returns name of this color if its available.
func (c Color) Name() (string, bool) {
	r, g, b, _ := c.RGBA255()
	rgb := [3]uint8{r, g, b}
	for k, v := range namedColors {
		if v == rgb {
			return k, true
		}
	}
	return "", false
}

// Implement the Go TextUnmarshaler interface
func (c *Color) UnmarshalText(text []byte) error {
	col, err := Parse(string(text))
	if err != nil {
		return err
	}
	c.R = col.R
	c.G = col.G
	c.B = col.B
	c.A = col.A
	return nil
}

// Implement the Go TextMarshaler interface
func (c Color) MarshalText() ([]byte, error) {
	return []byte(c.HexString()), nil
}

// FromHsv creates a Color from HSV colors.
//
// Arguments:
//
//   - h: Hue angle [0..360]
//   - s: Saturation [0..1]
//   - v: Value [0..1]
//   - a: Alpha [0..1]
func FromHsv(h, s, v, a float64) Color {
	r, g, b := hsvToRgb(normalizeAngle(h), clamp0_1(s), clamp0_1(v))
	return Color{r, g, b, clamp0_1(a)}
}

// FromHsl creates a Color from HSL colors.
//
// Arguments:
//
//   - h: Hue angle [0..360]
//   - s: Saturation [0..1]
//   - l: Lightness [0..1]
//   - a: Alpha [0..1]
func FromHsl(h, s, l, a float64) Color {
	r, g, b := hslToRgb(normalizeAngle(h), clamp0_1(s), clamp0_1(l))
	return Color{r, g, b, clamp0_1(a)}
}

// FromHwb creates a Color from HWB colors.
//
// Arguments:
//
//   - h: Hue angle [0..360]
//   - w: Whiteness [0..1]
//   - b: Blackness [0..1]
//   - a: Alpha [0..1]
func FromHwb(h, w, b, a float64) Color {
	r, g, b := hwbToRgb(normalizeAngle(h), clamp0_1(w), clamp0_1(b))
	return Color{r, g, b, clamp0_1(a)}
}

func fromLinear(x float64) float64 {
	if x >= 0.0031308 {
		return 1.055*math.Pow(x, 1.0/2.4) - 0.055
	}
	return 12.92 * x
}

// FromLinearRGB creates a Color from linear-light RGB colors.
//
// Arguments:
//
//   - r: Red value [0..1]
//   - g: Green value [0..1]
//   - b: Blue value [0..1]
//   - a: Alpha value [0..1]
func FromLinearRGB(r, g, b, a float64) Color {
	return Color{fromLinear(r), fromLinear(g), fromLinear(b), clamp0_1(a)}
}

// FromOklab creates a Color from Oklab colors.
//
// Arguments:
//
//   - l: Perceived lightness
//   - a: How green/red the color is
//   - b: How blue/yellow the color is
//   - alpha: Alpha [0..1]
func FromOklab(l, a, b, alpha float64) Color {
	l_ := math.Pow(l+0.3963377774*a+0.2158037573*b, 3)
	m_ := math.Pow(l-0.1055613458*a-0.0638541728*b, 3)
	s_ := math.Pow(l-0.0894841775*a-1.2914855480*b, 3)

	R := 4.0767416621*l_ - 3.3077115913*m_ + 0.2309699292*s_
	G := -1.2684380046*l_ + 2.6097574011*m_ - 0.3413193965*s_
	B := -0.0041960863*l_ - 0.7034186147*m_ + 1.7076147010*s_

	return FromLinearRGB(R, G, B, alpha)
}

// FromOklch creates a Color from OKLCh colors.
//
// Arguments:
//
//   - l: Perceived lightness
//   - c: Chroma
//   - h: Hue angle in radians
//   - alpha: Alpha [0..1]
func FromOklch(l, c, h, alpha float64) Color {
	return FromOklab(l, c*math.Cos(h), c*math.Sin(h), alpha)
}

var black = Color{0, 0, 0, 1}

// Parse parses CSS color string and returns, if successful, a Color.
func Parse(s string) (Color, error) {
	input := s
	s = strings.TrimSpace(strings.ToLower(s))

	if s == "transparent" {
		return Color{0, 0, 0, 0}, nil
	}

	// Predefined name / keyword
	c, ok := namedColors[s]
	if ok {
		return Color{float64(c[0]) / 255, float64(c[1]) / 255, float64(c[2]) / 255, 1}, nil
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
		s = s[op+1 : len(s)-1]
		f := func(c rune) bool {
			return c == ',' || c == '/' || c == ' '
		}
		params := strings.FieldsFunc(s, f)

		if len(params) != 3 && len(params) != 4 {
			return black, fmt.Errorf("Invalid format")
		}

		alpha := 1.0
		if len(params) == 4 {
			v, ok, _ := parsePercentOrFloat(params[3])
			if !ok {
				return black, fmt.Errorf("Invalid format")
			}
			alpha = clamp0_1(v)
		}

		if fname == "rgb" || fname == "rgba" {
			r, okR, _ := parsePercentOr255(params[0])
			g, okG, _ := parsePercentOr255(params[1])
			b, okB, _ := parsePercentOr255(params[2])

			if okR && okG && okB {
				return Color{
					clamp0_1(r),
					clamp0_1(g),
					clamp0_1(b),
					alpha,
				}, nil
			}
			return black, fmt.Errorf("Wrong %s() components, %s", fname, input)

		} else if fname == "hsl" || fname == "hsla" {
			h, okH := parseAngle(params[0])
			s, okS, _ := parsePercentOrFloat(params[1])
			l, okL, _ := parsePercentOrFloat(params[2])

			if okH && okS && okL {
				return FromHsl(h, s, l, alpha), nil
			}
			return black, fmt.Errorf("Wrong %s() components, %s", fname, input)

		} else if fname == "hwb" || fname == "hwba" {
			H, okH := parseAngle(params[0])
			W, okW, _ := parsePercentOrFloat(params[1])
			B, okB, _ := parsePercentOrFloat(params[2])

			if okH && okW && okB {
				return FromHwb(H, W, B, alpha), nil
			}
			return black, fmt.Errorf("Wrong hwb() components, %s", input)

		} else if fname == "hsv" || fname == "hsva" {
			h, okH := parseAngle(params[0])
			s, okS, _ := parsePercentOrFloat(params[1])
			v, okV, _ := parsePercentOrFloat(params[2])

			if okH && okS && okV {
				return FromHsv(h, s, v, alpha), nil
			}
			return black, fmt.Errorf("Wrong hsv() components, %s", input)

		} else if fname == "oklab" {
			l, okL, _ := parsePercentOrFloat(params[0])
			a, okA, fmtA := parsePercentOrFloat(params[1])
			b, okB, fmtB := parsePercentOrFloat(params[2])

			if okL && okA && okB {
				if fmtA {
					a = remap(a, -1.0, 1.0, -0.4, 0.4)
				}
				if fmtB {
					b = remap(b, -1.0, 1.0, -0.4, 0.4)
				}
				return FromOklab(math.Max(l, 0), a, b, alpha), nil
			}
			return black, fmt.Errorf("Wrong oklab() components, %s", input)

		} else if fname == "oklch" {
			l, okL, _ := parsePercentOrFloat(params[0])
			c, okC, fmtC := parsePercentOrFloat(params[1])
			h, okH := parseAngle(params[2])

			if okL && okC && okH {
				if fmtC {
					c = c * 0.4
				}
				return FromOklch(math.Max(l, 0), math.Max(c, 0), h*math.Pi/180, alpha), nil
			}
			return black, fmt.Errorf("Wrong oklch() components, %s", input)
		}
	}

	// RGB hexadecimal format without '#' prefix
	c2, ok2 := parseHex(s)
	if ok2 {
		return c2, nil
	}

	return black, fmt.Errorf("Invalid color format, %s", input)
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

func modulo(x, y float64) float64 {
	return math.Mod(math.Mod(x, y)+y, y)
}

func hueToRgb(n1, n2, h float64) float64 {
	h = modulo(h, 6)
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

// Returns (result, ok?, percentage?)
func parsePercentOrFloat(s string) (float64, bool, bool) {
	if strings.HasSuffix(s, "%") {
		f, ok := parseFloat(s[:len(s)-1])
		if ok {
			return f / 100, true, true
		}
		return 0, false, true
	}
	f, ok := parseFloat(s)
	return f, ok, false
}

// Returns (result, ok?, percentage?)
func parsePercentOr255(s string) (float64, bool, bool) {
	if strings.HasSuffix(s, "%") {
		f, ok := parseFloat(s[:len(s)-1])
		if ok {
			return f / 100, true, true
		}
		return 0, false, true
	}
	f, ok := parseFloat(s)
	if ok {
		return f / 255, true, false
	}
	return 0, false, false
}

// Result angle in degrees (not normalized)
func parseAngle(s string) (float64, bool) {
	if strings.HasSuffix(s, "deg") {
		return parseFloat(s[:len(s)-3])
	}
	if strings.HasSuffix(s, "grad") {
		f, ok := parseFloat(s[:len(s)-4])
		if ok {
			return f / 400 * 360, true
		}
		return 0, false
	}
	if strings.HasSuffix(s, "rad") {
		f, ok := parseFloat(s[:len(s)-3])
		if ok {
			return f / math.Pi * 180, true
		}
		return 0, false
	}
	if strings.HasSuffix(s, "turn") {
		f, ok := parseFloat(s[:len(s)-4])
		if ok {
			return f * 360, true
		}
		return 0, false
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

// Map t which is in range [a, b] to range [c, d]
func remap(t, a, b, c, d float64) float64 {
	return (t-a)*((d-c)/(b-a)) + c
}
