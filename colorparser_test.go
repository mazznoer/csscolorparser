package csscolorparser

import (
	"image/color"
	"testing"
)

// --- Helper functions

func test(t *testing.T, a, b interface{}) {
	if a != b {
		t.Helper()
		t.Errorf("left: %v, right: %v", a, b)
	}
}

func testTrue(t *testing.T, b bool) {
	if !b {
		t.Helper()
		t.Errorf("it false")
	}
}

func testColor(t *testing.T, a, b Color) {
	x := arr8(a.RGBA255())
	y := arr8(b.RGBA255())
	if x == y {
		return
	}
	t.Helper()
	t.Errorf("left: %v, right: %v", x, y)
}

func testGoColor(t *testing.T, a, b color.Color) {
	x := arr32(a.RGBA())
	y := arr32(b.RGBA())
	if x == y {
		return
	}
	t.Helper()
	t.Errorf("left: %v, right: %v", x, y)
}

func arr8(r, g, b, a uint8) [4]uint8 {
	return [4]uint8{r, g, b, a}
}

func arr32(r, g, b, a uint32) [4]uint32 {
	return [4]uint32{r, g, b, a}
}

// ---

func Test_Color(t *testing.T) {
	var c Color

	c = Color{0, 0, 1, 1}
	test(t, c.HexString(), "#0000ff")
	test(t, c.RGBString(), "rgb(0,0,255)")
	testGoColor(t, c, color.NRGBA{0, 0, 255, 255})

	c = Color{0, 0, 1, 0.5}
	test(t, c.HexString(), "#0000ff80")
	test(t, c.RGBString(), "rgba(0,0,255,0.5)")
	//testGoColor(t, c, color.NRGBA{0,0,255,127})

	testGoColor(t, Color{A: 1}, color.Gray{0})

	c = Color{1.2001, 0.999, -0.001, 0.001}.Clamp()
	testColor(t, c, Color{1, 0.999, 0, 0.001})

	c = FromHwb(0, 0, 0, 1)
	test(t, c.HexString(), "#ff0000")

	c = FromHwb(360, 0, 0, 1)
	test(t, c.HexString(), "#ff0000")

	c = FromHsv(120, 1, 1, 1)
	test(t, c.HexString(), "#00ff00")

	c = FromHsl(180, 1, 0.5, 1)
	test(t, c.HexString(), "#00ffff")

	c = FromOklab(0.62796, 0.22486, 0.12585, 1)
	test(t, c.HexString(), "#ff0000")

	c = FromOklch(0.62796, 0.25768, 0.51, 1)
	test(t, c.HexString(), "#ff0000")

	c = FromOklch(0.86644, 0.29483, 2.487, 1)
	test(t, c.HexString(), "#00ff00")
}

func Test_ParseColor(t *testing.T) {
	data0 := []struct {
		s     string
		rgba8 [4]uint8
	}{
		{"transparent", [4]uint8{0, 0, 0, 0}},
		{"rebeccapurple", [4]uint8{102, 51, 153, 255}},
		{"#ff00ff64", [4]uint8{255, 0, 255, 100}},
		{"ff00ff64", [4]uint8{255, 0, 255, 100}},
		{"rgb(247,179,99)", [4]uint8{247, 179, 99, 255}},
		{"rgb(50% 50% 50%)", [4]uint8{128, 128, 128, 255}},
		{"rgb(247,179,99,0.37)", [4]uint8{247, 179, 99, 94}},
		{"oklab(64.3% 52.6% 40% 2.5%)", [4]uint8{6, 26, 133, 6}},
		{"oklch(0.46212, 80.9%, 29.23388, 17.33713)", [4]uint8{214, 206, 150, 255}},
		{"hsl(270 0% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hwb(0 50% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hsv(0 0% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hsv(0 0% 100%)", [4]uint8{255, 255, 255, 255}},
		{"hsv(0 0% 19%)", [4]uint8{48, 48, 48, 255}},
		//{"lab(0%,0,0)", [4]uint8{0, 0, 0, 255}},
		//{"lab(100%,0,0)", [4]uint8{255, 255, 255, 255}},
		//{"lch(0%,0,0)", [4]uint8{0, 0, 0, 255}},
		//{"lch(100%,0,0)", [4]uint8{255, 255, 255, 255}},
	}
	for _, d := range data0 {
		c, err := Parse(d.s)
		test(t, err, nil)
		test(t, arr8(c.RGBA255()), d.rgba8)
	}

	data1 := []struct {
		s string
		c Color
	}{
		{"hwb(0, 0%, 0%)", FromHwb(0, 0, 0, 1)},
		{"hwb(320, 10%, 30%)", FromHwb(320, 0.1, 0.3, 1)},
		{"hsv(120, 30%, 50%)", FromHsv(120, 0.3, 0.5, 1)},
		{"hsl(120, 30%, 50%)", FromHsl(120, 0.3, 0.5, 1)},
	}
	for _, dt := range data1 {
		c, err := Parse(dt.s)
		test(t, err, nil)
		test(t, dt.c.HexString(), c.HexString())
	}

	data2 := []string{
		"#666666",
		"#ff0000",
		"#00ff7f",
	}
	for _, s := range data2 {
		c1, err := Parse(s)
		test(t, err, nil)
		c2, err2 := Parse(c1.RGBString())
		test(t, err2, nil)
		test(t, c2.HexString(), s)
	}

	/*
	a, err := Parse("#7654CD")
	test(t, err, nil)
	b, err := Parse("lab(44.36% 36.05 -58.99)")
	test(t, err, nil)
	testColor(t, a, b)
	*/
}

func Test_MarshalUnmarshal(t *testing.T) {
	var c Color
	err := c.UnmarshalText([]byte("gold"))
	test(t, err, nil)
	test(t, c.HexString(), "#ffd700")

	encoding, err := c.MarshalText()
	test(t, err, nil)
	test(t, string(encoding), "#ffd700")

	err = c.UnmarshalText([]byte("golden"))
	testTrue(t, err != nil)
}

func Test_EqualColorsBlack(t *testing.T) {
	data := []string{
		"black",
		"#000",
		"#000f",
		"#000000",
		"#000000ff",
		"rgb(0,0,0)",
		"rgb(0% 0% 0%)",
		"rgb(0 0 0 100%)",
		"hsl(270,100%,0%)",
		"hwb(90 0% 100%)",
		"hwb(120deg 0% 100% 100%)",
		"hsv(120 100% 0%)",
	}
	for _, d := range data {
		c, err := Parse(d)
		test(t, err, nil)
		testColor(t, c, Color{0, 0, 0, 1})
	}
}

func Test_EqualColorsRed(t *testing.T) {
	data := []string{
		"red",
		"#f00",
		"#f00f",
		"#ff0000",
		"#ff0000ff",
		"rgb(255,0,0)",
		"rgb(255 0 0)",
		"rgb(700, -99, 0)", // clamp to 0..255
		"rgb(100% 0% 0%)",
		"rgb(200% -10% -100%)", // clamp to 0%..100%
		"rgb(255 0 0 100%)",
		"RGB( 255 , 0 , 0 )",
		"RGB( 255   0   0 )",
		"hsl(0,100%,50%)",
		"hsl(360 100% 50%)",
		"hwb(0 0% 0%)",
		"hwb(360deg 0% 0% 100%)",
		"hsv(0 100% 100%)",
		"oklab(0.62796, 0.22486, 0.12585)",
		"oklch(0.62796, 0.25768, 29.23388)",
	}
	for _, d := range data {
		c, err := Parse(d)
		test(t, err, nil)
		testColor(t, c, Color{1, 0, 0, 1})
	}
}

func Test_EqualColorsLime(t *testing.T) {
	data := []string{
		"lime",
		"#0f0",
		"#0f0f",
		"#00ff00",
		"#00ff00ff",
		"rgb(0,255,0)",
		"rgb(0% 100% 0%)",
		"rgb(0 255 0 / 100%)",
		"rgba(0,255,0,1)",
		"hsl(120,100%,50%)",
		"hsl(120deg 100% 50%)",
		"hsl(-240 100% 50%)",
		"hsl(-240deg 100% 50%)",
		"hsl(0.3333turn 100% 50%)",
		"hsl(133.333grad 100% 50%)",
		"hsl(2.0944rad 100% 50%)",
		"hsla(120,100%,50%,100%)",
		"hwb(120 0% 0%)",
		"hwb(480deg 0% 0% / 100%)",
		"hsv(120 100% 100%)",
		"oklab(0.86644, -0.23389, 0.1795)",
		"oklch(0.86644, 0.29483, 142.49535)",
	}
	for _, d := range data {
		c, err := Parse(d)
		test(t, err, nil)
		testColor(t, c, Color{0, 1, 0, 1})
	}
}

func Test_EqualColorsLimeAlpha(t *testing.T) {
	data := []string{
		"#00ff0080",
		"rgb(0,255,0,50%)",
		"rgb(0% 100% 0% / 0.5)",
		"rgba(0%,100%,0%,50%)",
		"hsl(120,100%,50%,0.5)",
		"hsl(120deg 100% 50% / 50%)",
		"hsla(120,100%,50%,0.5)",
		"hwb(120 0% 0% / 50%)",
		"hsv(120 100% 100% / 50%)",
	}
	for _, d := range data {
		c, err := Parse(d)
		test(t, err, nil)
		testColor(t, c, Color{0, 1, 0, 0.5})
	}
}

func Test_InvalidData(t *testing.T) {
	data := []string{
		"",
		"bloodred",
		"#78afzd",
		"#fffff",
		"rgb(0,255,8s)",
		"rgb(100%,z9%,75%)",
		"cmyk(1 0 0)",
		"rgba(0 0)",
		"hsl(90',100%,50%)",
		"hsl(deg 100% 50%)",
		"hsl(Xturn 100% 50%)",
		"hsl(Zgrad 100% 50%)",
		"hsl(180 1 x%)",
		"hsla(360)",
		"hwb(Xrad,50%,50%)",
		"hwb(270 0% 0% 0% 0%)",
		"hsv(120 100% 100% 1 50%)",
		"hsv(120 XXX 100%)",
		"oklab(0,0)",
		"oklab(0,0,x,0)",
		"oklch(0,0,0,0,0)",
		"oklch(0,0,0,x)",
	}
	for _, s := range data {
		c, err := Parse(s)
		testTrue(t, err != nil)
		testColor(t, c, Color{A: 1})
	}
}

func Test_Utils(t *testing.T) {
	// parseAngle

	data := []struct {
		s string
		f float64
	}{
		{"360", 360},
		{"127.356", 127.356},
		{"+120deg", 120},
		{"90deg", 90},
		{"-127deg", -127},
		{"100grad", 90},
		{"1.5707963267948966rad", 90},
		{"0.25turn", 90},
		{"-0.25turn", -90},
	}
	for _, s := range data {
		d, ok := parseAngle(s.s)
		testTrue(t, ok)
		test(t, d, s.f)
	}

	// normalizeAngle

	data2 := [][2]float64{
		{0, 0},
		{360, 0},
		{400, 40},
		{1155, 75},
		{-360, 0},
		{-90, 270},
		{-765, 315},
	}
	for _, d := range data2 {
		test(t, normalizeAngle(d[0]), d[1])
	}
}
