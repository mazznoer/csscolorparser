package csscolorparser

import (
	"fmt"
	"math"
)

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
	a := uint8(clamp(c.A) * 100)
	if a < 100 {
		return fmt.Sprintf("rgb(%d %d %d / %d%%)", r, g, b, a)
	}
	return fmt.Sprintf("rgb(%d %d %d)", r, g, b)
}

// Name returns name of this color if its available.
func (c Color) Name() (string, bool) {
	r, g, b, _ := c.RGBA255()
	rgb := [3]uint8{r, g, b}
	for k, v := range NamedColors {
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
	r, g, b := hsvToRgb(normalizeAngle(h), clamp(s), clamp(v))
	return Color{r, g, b, clamp(a)}
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
	r, g, b := hslToRgb(normalizeAngle(h), clamp(s), clamp(l))
	return Color{r, g, b, clamp(a)}
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
	r, g, b := hwbToRgb(normalizeAngle(h), clamp(w), clamp(b))
	return Color{r, g, b, clamp(a)}
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
	return Color{fromLinear(r), fromLinear(g), fromLinear(b), clamp(a)}
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

func FromLab(l, a, b, alpha float64) Color {
	x, y, z := labToXyz(l, a, b)
	R, G, B := xyzToLinearRgb(x, y, z)
	return FromLinearRGB(R, G, B, alpha)
}

func FromLch(l, c, h, alpha float64) Color {
	return FromLab(l, c*math.Cos(h), c*math.Sin(h), alpha)
}

// Utils

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
	r = clamp(hueToRgb(n1, n2, h+2))
	g = clamp(hueToRgb(n1, n2, h))
	b = clamp(hueToRgb(n1, n2, h-2))
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

func labToXyz(l, a, b float64) (x, y, z float64) {
	const (
		e  = 216.0 / 24389.0
		k  = 24389.0 / 27.0
		Xn = 0.95047
		Yn = 1.00000
		Zn = 1.08883
	)

	fy := (l + 16.0) / 116.0
	fx := fy + (a / 500.0)
	fz := fy - (b / 200.0)

	// Calculate Y
	if l > 8.0 { // Standard check is usually on L rather than fy^3 for Y
		y = Yn * math.Pow(fy, 3)
	} else {
		y = Yn * (l / k)
	}

	// Calculate X
	fx3 := math.Pow(fx, 3)
	if fx3 > e {
		x = Xn * fx3
	} else {
		x = Xn * (116.0*fx - 16.0) / k
	}

	// Calculate Z
	fz3 := math.Pow(fz, 3)
	if fz3 > e {
		z = Zn * fz3
	} else {
		z = Zn * (116.0*fz - 16.0) / k
	}

	return
}

func xyzToLinearRgb(x, y, z float64) (r, g, b float64) {
	r = 3.2404537*x - 1.537177*y - 0.498531*z
	g = -0.969266*x + 1.876010*y + 0.041556*z
	b = 0.055643*x - 0.204025*y + 1.057225*z
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))
	return
}
