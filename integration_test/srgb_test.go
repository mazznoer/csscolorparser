package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestSrgbAbsolute(t *testing.T) {
	testData := [][2]string{
		// bare number
		{"#a83935", "color(srgb 0.658 0.223 0.206)"},
		{"#d33a48", "color(srgb 0.829 0.227 0.282)"},
		{"#a2e65f", "color(srgb 0.637 0.903 0.372)"},
		{"#5d01ae", "color(srgb 0.364 0.003 0.684)"},
		{"#7104e8", "color(srgb 0.442 0.017 0.911)"},
		{"#c57527", "color(srgb 0.774 0.459 0.151)"},
		{"#ccb082", "color(srgb 0.801 0.69 0.508)"},
		{"#5f1650", "color(srgb 0.373 0.087 0.315)"},
		{"#febbf7", "color(srgb 0.996 0.733 0.97)"},
		{"#1c6ed2", "color(srgb 0.108 0.431 0.824)"},
		{"#6f8404", "color(srgb 0.436 0.519 0.015)"},
		{"#f1c7a7", "color(srgb 0.946 0.779 0.656)"},
		{"#1daa7a", "color(srgb 0.113 0.668 0.477)"},
		{"#e24fff", "color(srgb 0.886 0.31 0.999)"},
		{"#036a14", "color(srgb 0.012 0.415 0.077)"},
		// percentage
		{"#478c08", "color(srgb 28% 55% 3%)"},
		{"#177d7d", "color(srgb 9% 49% 49%)"},
		{"#ad30a6", "color(srgb 68% 19% 65%)"},
		{"#52a1a3", "color(srgb 32% 63% 64%)"},
		{"#9eadf0", "color(srgb 62% 68% 94%)"},
		{"#bfa60a", "color(srgb 75% 65% 4%)"},
		{"#128fcf", "color(srgb 7% 56% 81%)"},
		{"#2eedc4", "color(srgb 18% 93% 77%)"},
		{"#f72170", "color(srgb 97% 13% 44%)"},
		{"#420a2e", "color(srgb 26% 4% 18%)"},
		{"#d9f2d1", "color(srgb 85% 95% 82%)"},
		{"#004d6e", "color(srgb 0% 30% 43%)"},
		{"#99ed91", "color(srgb 60% 93% 57%)"},
		{"#b5bf5e", "color(srgb 71% 75% 37%)"},
		{"#4dd9c2", "color(srgb 30% 85% 76%)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
