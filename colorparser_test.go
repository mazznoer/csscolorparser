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
	testData := []colorPair{
		{"transparent", color.RGBA{0, 0, 0, 0}},
		{"rebeccapurple", color.RGBA{102, 51, 153, 255}},
		{"#ff00ff64", color.NRGBA{255, 0, 255, 100}},
		{"rgb(247,179,99)", color.NRGBA{247, 179, 99, 255}},
		{"rgb(50% 50% 50%)", color.NRGBA{127, 127, 127, 255}},
		{"rgb(247,179,99,0.37)", color.NRGBA{247, 179, 99, 94}},
		{"hsl(270 0% 50%)", color.NRGBA{127, 127, 127, 255}},
		{"hwb(0 50% 50%)", color.NRGBA{127, 127, 127, 255}},
		{"hsv(0 0% 50%)", color.NRGBA{127, 127, 127, 255}},
		{"hsv(0 0% 100%)", color.NRGBA{255, 255, 255, 255}},
		{"hsv(0 0% 19%)", color.NRGBA{48, 48, 48, 255}},
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
		"rgb(0% 100% 0)",
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
	lime := color.NRGBA{0, 255, 0, 255}
	for _, d := range data {
		c, _ := Parse(d)
		if !isColorEqual(lime, c) {
			t.Errorf("Not lime, %s -> %v", d, c)
			break
		}
	}
}

func TestEqualColorsLimeAlpha(t *testing.T) {
	data := []string{
		"#00ff007f",
		"rgb(0,255,0,50%)",
		"rgb(0% 100% 0% / 0.5)",
		"rgba(0%,100%,0%,50%)",
		"hsl(120,100%,50%,0.5)",
		"hsl(120deg 100% 50% / 50%)",
		"hsla(120,100%,50%,0.5)",
		"hwb(120 0% 0% / 50%)",
		"hsv(120 100% 100% / 50%)",
	}
	limeAlpha := color.NRGBA{0, 255, 0, 127}
	for _, d := range data {
		c, _ := Parse(d)
		if !isColorEqual(limeAlpha, c) {
			t.Errorf("Not lime 0.5 alpha, %s -> %v", d, c)
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
		//"rgb (127,255,0)",
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
