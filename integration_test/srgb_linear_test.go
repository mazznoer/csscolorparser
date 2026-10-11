package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestSrgbLinearAbsolute(t *testing.T) {
	testData := [][2]string{
		// bare number
		{"#f578e8", "color(srgb-linear 0.91 0.189 0.807)"},
		{"#d1f483", "color(srgb-linear 0.636 0.907 0.227)"},
		{"#4f1f91", "color(srgb-linear 0.079 0.014 0.285)"},
		{"#fa8d9b", "color(srgb-linear 0.955 0.265 0.327)"},
		{"#a710d6", "color(srgb-linear 0.386 0.005 0.674)"},
		{"#c2e167", "color(srgb-linear 0.538 0.752 0.136)"},
		{"#a4b38d", "color(srgb-linear 0.37 0.452 0.265)"},
		{"#778b81", "color(srgb-linear 0.186 0.26 0.221)"},
		{"#cea43e", "color(srgb-linear 0.618 0.371 0.048)"},
		{"#4db369", "color(srgb-linear 0.075 0.451 0.142)"},
		{"#d8e591", "color(srgb-linear 0.686 0.785 0.282)"},
		{"#ede5f2", "color(srgb-linear 0.845 0.783 0.884)"},
		{"#3b718a", "color(srgb-linear 0.044 0.164 0.254)"},
		{"#51bddd", "color(srgb-linear 0.082 0.511 0.725)"},
		{"#e8e7c3", "color(srgb-linear 0.804 0.802 0.544)"},
		// percentage
		{"#da27b1", "color(srgb-linear 70% 2% 44%)"},
		{"#f6e2f1", "color(srgb-linear 92% 76% 88%)"},
		{"#efbcfd", "color(srgb-linear 86% 50% 98%)"},
		{"#6f61ba", "color(srgb-linear 16% 12% 49%)"},
		{"#6c866f", "color(srgb-linear 15% 24% 16%)"},
		{"#dbd3c0", "color(srgb-linear 71% 65% 53%)"},
		{"#3f00ce", "color(srgb-linear 5% 0% 62%)"},
		{"#a2ce93", "color(srgb-linear 36% 62% 29%)"},
		{"#50b6af", "color(srgb-linear 8% 47% 43%)"},
		{"#da5576", "color(srgb-linear 70% 9% 18%)"},
		{"#89d181", "color(srgb-linear 25% 64% 22%)"},
		{"#cdfa6c", "color(srgb-linear 61% 96% 15%)"},
		{"#b19bd0", "color(srgb-linear 44% 33% 63%)"},
		{"#65d7ec", "color(srgb-linear 13% 68% 84%)"},
		{"#a0a0b8", "color(srgb-linear 35% 35% 48%)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
