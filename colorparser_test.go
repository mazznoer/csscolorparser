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

func TestParseColor(t *testing.T) {
	type colorPair struct {
		in  string
		out color.Color
	}

	red := color.RGBA{255, 0, 0, 255}
	redAlpha := color.RGBA{255, 0, 0, 127}

	testData := []colorPair{
		{"transparent", color.RGBA{0, 0, 0, 0}},
		{"red", red},
		{"#f00", red},
		{"#ff0000", red},
		{"#f00f", red},
		{"#ff0000ff", red},
		//{"#ff000065", color.RGBA{255, 0, 0, 101}},
		{"rgb(255,0,0)", red},
		{"rgb(255 0 0)", red},
		{"RGB( 255 , 0 , 0 )", red},
		{"RGB( 255   0   0 )", red},
		{"rgb(255 0 0 / 50%)", redAlpha},
		{"rgb(100% 0% 0% / 0.5)", redAlpha},
		{"rgba(255,0,0,0.5)", redAlpha},
		{"rgba(255,0,0,50%)", redAlpha},
		{"rgb(100%,0%,0%)", red},
		{"rgb(100% 0% 0%)", red},
		{"rgba(100%,0%,0%,0.5)", redAlpha},
		{"rgba(100%,0%,0%,50%)", redAlpha},
		{"hsl(0,100%,50%)", red},
		{"hsl(360 100% 50%)", red},
		{"hsl(0 100% 50% / 50%)", redAlpha},
		{"hsla(0,100%,50%,0.5)", redAlpha},
		{"hsl(360deg,100%,50%)", red},
		{"hsl(400grad,100%,50%)", red},
		{"hsl(0rad,100%,50%)", red},
		{"hsl(1turn,100%,50%)", red},
		{"hsl(4turn,100%,50%)", red},
		{"HSL(270 0% 50%)", color.RGBA{127, 127, 127, 255}},
		{"hwb(0 0% 0%)", red},
		{"hwb(0 0% 0% 50%)", redAlpha},
	}

	for _, d := range testData {
		c, err := Parse(d.in)
		if err != nil {
			t.Errorf("Parse error: %s", d.in)
			continue
		}
		if !isColorEqual(c, d.out) {
			t.Errorf("%s -> %v != %v", d.in, d.out, c)
		}
	}
}

func TestEqualColors1(t *testing.T) {
	data := []string{
		"black",
		"#000",
		"#000000",
		"rgb(0,0,0)",
		"rgb(0% 0% 0%)",
		"rgb(0 0 0 100%)",
		"hsl(270,100%,0%)",
		"hwb(90 0% 100%)",
		"hwb(120deg 0% 100% 100%)",
	}
	c, _ := Parse(data[0])
	t.Log(c)
	for _, d := range data {
		cx, _ := Parse(d)
		if !isColorEqual(cx, c) {
			t.Errorf("Not equal, %v -> %v", d, cx)
			break
		}
	}
}

func TestEqualColors2(t *testing.T) {
	data := []string{
		"lime",
		"#0f0",
		"#0f0f",
		"#00ff00",
		"#00ff00ff",
		"rgb(0,255,0)",
		"rgb(0% 100% 0)",
		"hsl(120,100%,50%)",
		"hsl(120deg 100% 50%)",
		"hsl(-240 100% 50%)",
		"hsl(-240deg 100% 50%)",
		"hsl(0.3333turn 100% 50%)",
		"hsl(133.333grad 100% 50%)",
		"hsl(2.0944rad 100% 50%)",
		"hwb(120 0% 0%)",
	}
	c, _ := Parse(data[0])
	t.Log(c)
	for _, d := range data {
		cx, _ := Parse(d)
		if !isColorEqual(cx, c) {
			t.Errorf("Not equal, %v", d)
			break
		}
	}
}

func TestEqualColors3(t *testing.T) {
	data := []string{
		"hsl(90,100%,50%)",
		"hsl(-270deg 100% 50% 1)",
		"hsl(100grad 100% 50%)",
		"hsl(-0.75turn 1 0.5)",
		"hsla(1.25turn 100% 50% 100%)",
		"hwb(450 0% 0%)",
		"hwb(0.25turn 0% 0%)",
		"hwb(90deg 0% 0% 100%)",
	}
	c, _ := Parse(data[0])
	t.Log(c)
	for _, d := range data {
		cx, _ := Parse(d)
		if !isColorEqual(cx, c) {
			t.Errorf("Not equal, %v", d)
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
		"rgb(255,0 0)",
		"rgb(0,255,8s)",
		"rgb(100%,z9%,75%)",
		"rgb (127,255,0)",
		"rgba(0 0)",
		"hsl(90degs,100%,50%)",
		"hsl(180 1 x%)",
		"hsla(360)",
		"hwb(Zdeg,50%,50%)",
	}
	for _, d := range testData {
		c, err := Parse(d)
		if err == nil {
			t.Errorf("It should fail, %s -> %v", d, c)
		}
		t.Log(err)
	}
}
