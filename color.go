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
	A := clamp(c.A)
	r = uint32(clamp(c.R)*A*65535 + 0.5)
	g = uint32(clamp(c.G)*A*65535 + 0.5)
	b = uint32(clamp(c.B)*A*65535 + 0.5)
	a = uint32(A*65535 + 0.5)
	return
}

// RGBA255 returns R, G, B, A values in the range 0..255
func (c Color) RGBA255() (r, g, b, a uint8) {
	r = uint8(clamp(c.R)*255 + 0.5)
	g = uint8(clamp(c.G)*255 + 0.5)
	b = uint8(clamp(c.B)*255 + 0.5)
	a = uint8(clamp(c.A)*255 + 0.5)
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
	R, G, B := labToLinearRGB(l, a, b)
	return FromLinearRGB(R, G, B, alpha)
}

func FromLch(l, c, h, alpha float64) Color {
	if math.IsNaN(h) {
		h = 0
	}
	return FromLab(l, c*math.Cos(h), c*math.Sin(h), alpha)
}

func FromXyzD50(x, y, z, alpha float64) Color {
	r, g, b := xyzD50ToLinearRGB(x, y, z)
	return FromLinearRGB(r, g, b, alpha)
}

func FromXyzD65(x, y, z, alpha float64) Color {
	r, g, b := xyzD65ToLinearRgb(x, y, z)
	return FromLinearRGB(r, g, b, alpha)
}

func (c Color) Hsl() [4]float64 {
	c = c.Clamp()
	h, s, l := rgbToHsl(c.R, c.G, c.B)
	return [4]float64{h, s, l, c.A}
}

func (c Color) Hwb() [4]float64 {
	c = c.Clamp()
	h, w, b := rgbToHwb(c.R, c.G, c.B)
	return [4]float64{h, w, b, c.A}
}

func (c Color) LinearRgb() [4]float64 {
	toLinear := func(x float64) float64 {
		if x >= 0.04045 {
			return math.Pow((x+0.055)/1.055, 2.4)
		}
		return x / 12.92
	}
	return [4]float64{
		toLinear(c.R),
		toLinear(c.G),
		toLinear(c.B),
		c.A,
	}
}

func (c Color) Lab() [4]float64 {
	rgb := c.LinearRgb()
	l, a, b := linearRGBToLab(rgb[0], rgb[1], rgb[2])
	return [4]float64{l, a, b, c.A}
}

func (col Color) Lch() [4]float64 {
	lab := col.Lab()
	l := lab[0]
	a := lab[1]
	b := lab[2]
	c := math.Hypot(a, b)
	h := math.Atan2(b, a)
	if h < 0 {
		h += math.Pi * 2
	}
	return [4]float64{l, c, h, col.A}
}

func (col Color) Oklab() [4]float64 {
	arr := col.LinearRgb()
	l := math.Cbrt(0.4121656120*arr[0] + 0.5362752080*arr[1] + 0.0514575653*arr[2])
	m := math.Cbrt(0.2118591070*arr[0] + 0.6807189584*arr[1] + 0.1074065790*arr[2])
	s := math.Cbrt(0.0883097947*arr[0] + 0.2818474174*arr[1] + 0.6302613616*arr[2])
	return [4]float64{
		0.2104542553*l + 0.7936177850*m - 0.0040720468*s,
		1.9779984951*l - 2.4285922050*m + 0.4505937099*s,
		0.0259040371*l + 0.7827717662*m - 0.8086757660*s,
		col.A,
	}
}

func (col Color) Oklch() [4]float64 {
	lab := col.Oklab()
	a := lab[1]
	b := lab[2]
	c := math.Sqrt(a*a + b*b)
	h := math.Atan2(b, a)
	return [4]float64{lab[0], c, h, lab[3]}
}

func (c Color) XyzD50() [4]float64 {
	rgb := c.LinearRgb()
	x, y, z := linearRGBToXYZD50(rgb[0], rgb[1], rgb[2])
	return [4]float64{x, y, z, c.A}
}

func (c Color) XyzD65() [4]float64 {
	rgb := c.LinearRgb()
	x, y, z := linearRgbToXyzD65(rgb[0], rgb[1], rgb[2])
	return [4]float64{x, y, z, c.A}
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

// ---

// hue returns the hue in degrees [0, 360) for sRGB values in [0, 1],
// given the precomputed max and min channel values.
// Achromatic colors (max == min) return 0.
func hue(r, g, b, max, min float64) float64 {
	delta := max - min
	if delta == 0 {
		return 0
	}

	var h float64
	switch max {
	case r:
		h = math.Mod((g-b)/delta, 6)
	case g:
		h = (b-r)/delta + 2
	default:
		h = (r-g)/delta + 4
	}

	h *= 60
	if h < 0 {
		h += 360
	}
	return h
}

// rgbToHsl converts sRGB (r, g, b in [0, 1]) to HSL.
// h is in degrees [0, 360); s and l are in [0, 1].
func rgbToHsl(r, g, b float64) (h, s, l float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	l = (max + min) / 2
	if delta == 0 {
		return 0, 0, l
	}

	s = delta / (1 - math.Abs(2*l-1))
	h = hue(r, g, b, max, min)
	return h, s, l
}

// rgbToHwb converts sRGB (r, g, b in [0, 1]) to HWB.
// h is in degrees [0, 360); w (whiteness) and bl (blackness) are in [0, 1].
func rgbToHwb(r, g, b float64) (h, w, bl float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h = hue(r, g, b, max, min)
	w = min
	bl = 1 - max
	return h, w, bl
}

// converts linear sRGB (Rec. 709 primaries, D65 white point)
// to CIE XYZ (D65). Inputs are in [0, 1]; white (1,1,1) maps to Y = 1.
func linearRgbToXyzD65(r, g, b float64) (x, y, z float64) {
	x = 0.4123907992659595*r + 0.357584339383878*g + 0.1804807884018343*b
	y = 0.21263900587151036*r + 0.715168678767756*g + 0.07219231536073371*b
	z = 0.01933081871559185*r + 0.11919477979462598*g + 0.9505321522496606*b
	return
}

// converts CIE XYZ (D65) to linear sRGB.
// Results are not clamped, so out-of-gamut colors can fall outside [0, 1].
func xyzD65ToLinearRgb(x, y, z float64) (r, g, b float64) {
	r = 3.2409699419045226*x - 1.537383177570094*y - 0.4986107602930034*z
	g = -0.9692436362808796*x + 1.8759675015077202*y + 0.04155505740717559*z
	b = 0.05563007969699366*x - 0.20397695888897652*y + 1.0569715142428786*z
	return
}
