package csscolorparser

import (
	"image/color"
	"testing"
)

type rgba struct {
	r, g, b, a uint32
}

func isColorEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	a := rgba{r1, g1, b1, a1}
	b := rgba{r2, g2, b2, a2}
	if a == b {
		return true
	}
	return false
}

func equalStr(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("%s != %s", a, b)
	}
}

func TestColor(t *testing.T) {
	a := Color{0, 0, 1, 1}
	equalStr(t, a.HexString(), "#0000ff")
	equalStr(t, a.RGBString(), "rgb(0,0,255)")

	b := Color{0, 0, 1, 0.5}
	equalStr(t, b.HexString(), "#0000ff80")
	equalStr(t, b.RGBString(), "rgba(0,0,255,0.5)")
}

func TestParseColor(t *testing.T) {
	type colorPair struct {
		in  string
		out [4]uint8
	}
	testData := []colorPair{
		{"transparent", [4]uint8{0, 0, 0, 0}},
		{"rebeccapurple", [4]uint8{102, 51, 153, 255}},
		{"#ff00ff64", [4]uint8{255, 0, 255, 100}},
		{"ff00ff64", [4]uint8{255, 0, 255, 100}},
		{"rgb(247,179,99)", [4]uint8{247, 179, 99, 255}},
		{"rgb(50% 50% 50%)", [4]uint8{128, 128, 128, 255}},
		{"rgb(247,179,99,0.37)", [4]uint8{247, 179, 99, 94}},
		{"hsl(270 0% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hwb(0 50% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hsv(0 0% 50%)", [4]uint8{128, 128, 128, 255}},
		{"hsv(0 0% 100%)", [4]uint8{255, 255, 255, 255}},
		{"hsv(0 0% 19%)", [4]uint8{48, 48, 48, 255}},
	}
	for _, d := range testData {
		c, err := Parse(d.in)
		if err != nil {
			t.Errorf("Parse error: %s", d.in)
			continue
		}
		r, g, b, a := c.RGBA255()
		rgba := [4]uint8{r, g, b, a}
		if rgba != d.out {
			t.Errorf("%s -> %v != %v", d.in, d.out, rgba)
		}
	}
}

func TestNamedColors(t *testing.T) {
	data := [][2]string{
		{"aliceblue", "#f0f8ff"},
		{"bisque", "#ffe4c4"},
		{"chartreuse", "#7fff00"},
		{"coral", "#ff7f50"},
		{"crimson", "#dc143c"},
		{"dodgerblue", "#1e90ff"},
		{"firebrick", "#b22222"},
		{"gold", "#ffd700"},
		{"hotpink", "#ff69b4"},
		{"indigo", "#4b0082"},
		{"lavender", "#e6e6fa"},
		{"plum", "#dda0dd"},
		{"salmon", "#fa8072"},
		{"skyblue", "#87ceeb"},
		{"tomato", "#ff6347"},
		{"violet", "#ee82ee"},
		{"yellowgreen", "#9acd32"},
	}
	for _, d := range data {
		c, _ := Parse(d[0])
		hex := c.HexString()
		if hex != d[1] {
			t.Errorf("%s != %s", hex, d[1])
		}
	}
}

func TestEqualColorsBlack(t *testing.T) {
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
	black := color.NRGBA{0, 0, 0, 255}
	for _, d := range data {
		c, _ := Parse(d)
		if !isColorEqual(black, c) {
			t.Errorf("Not black, %s -> %v", d, c)
			break
		}
	}
}

func TestEqualColorsRed(t *testing.T) {
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
	}
	red := color.NRGBA{255, 0, 0, 255}
	for _, d := range data {
		c, _ := Parse(d)
		if !isColorEqual(red, c) {
			t.Errorf("Not red, %s -> %v", d, c)
			break
		}
	}
}

func TestEqualColorsLime(t *testing.T) {
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
	}
	lime := [4]uint8{0, 255, 0, 255}
	for _, d := range data {
		c, _ := Parse(d)
		r, g, b, a := c.RGBA255()
		rgba := [4]uint8{r, g, b, a}
		if rgba != lime {
			t.Errorf("Not lime, %s -> %v", d, rgba)
			break
		}
	}
}

func TestEqualColorsLimeAlpha(t *testing.T) {
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
	limeAlpha := [4]uint8{0, 255, 0, 128}
	for _, d := range data {
		c, _ := Parse(d)
		r, g, b, a := c.RGBA255()
		rgba := [4]uint8{r, g, b, a}
		if rgba != limeAlpha {
			t.Errorf("Not lime 0.5 alpha, %s -> %v", d, rgba)
			break
		}
	}
}

func TestInvalidData(t *testing.T) {
	testData := []string{
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
	}
	for _, d := range testData {
		c, err := Parse(d)
		if err == nil {
			t.Errorf("It should fail, %s -> %v", d, c)
		}
		t.Log(err)
	}
}

func TestParseAngle(t *testing.T) {
	type pair struct {
		in  string
		out float64
	}
	testData := []pair{
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
	for _, s := range testData {
		d, ok := parseAngle(s.in)
		if !ok {
			t.Errorf("Parse error, %s", s.in)
		}
		if d != s.out {
			t.Errorf("%s -> %v != %v", s.in, d, s.out)
		}
	}
}

func TestNormalizeAngle(t *testing.T) {
	testData := [][2]float64{
		{0, 0},
		{360, 0},
		{400, 40},
		{1155, 75},
		{-360, 0},
		{-90, 270},
		{-765, 315},
	}
	for _, s := range testData {
		d := normalizeAngle(s[0])
		if d != s[1] {
			t.Errorf("%v -> %v != %v", s[0], d, s[1])
		}
	}
}
