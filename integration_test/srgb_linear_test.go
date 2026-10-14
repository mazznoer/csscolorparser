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

func TestSrgbLinearRelative(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "color(from #bad455 srgb-linear none none none)"},
		{"#bad45500", "color(from #bad455 srgb-linear r g b / none)"},

		// bare number
		{"#d8b1d5", "color(from #fff srgb-linear 0.687 0.437 0.665)"},
		{"#73ff6f", "color(from #fff srgb-linear 0.173 0.997 0.159)"},
		{"#86badb", "color(from #fff srgb-linear 0.24 0.491 0.711)"},
		{"#21a2da", "color(from #fff srgb-linear 0.015 0.361 0.704)"},
		{"#f5d3e3", "color(from #fff srgb-linear 0.909 0.65 0.769)"},
		{"#c7e0f5", "color(from #fff srgb-linear 0.571 0.742 0.911)"},
		//{"#dcf4f8", "color(from #fff srgb-linear 0.712 0.906 0.935)"},
		{"#a9a3d0", "color(from #fff srgb-linear 0.399 0.364 0.63)"},
		{"#58af31", "color(from #fff srgb-linear 0.097 0.43 0.031)"},
		{"#c7c13a", "color(from #fff srgb-linear 0.57 0.531 0.043)"},
		{"#b1b0c9", "color(from #fff srgb-linear 0.439 0.435 0.585)"},
		{"#e0d3e4", "color(from #fff srgb-linear 0.744 0.653 0.776)"},
		{"#ed6bc1", "color(from #fff srgb-linear 0.844 0.147 0.531)"},
		{"#86daf8", "color(from #fff srgb-linear 0.238 0.7 0.942)"},
		{"#ebb56a", "color(from #fff srgb-linear 0.834 0.463 0.143)"},

		// percentage
		{"#559797", "color(from #fff srgb-linear 9% 31% 31%)"},
		{"#f2edba", "color(from #fff srgb-linear 89% 85% 49%)"},
		{"#005559", "color(from #fff srgb-linear 0% 9% 10%)"},
		{"#847eec", "color(from #fff srgb-linear 23% 21% 84%)"},
		{"#5dc7b1", "color(from #fff srgb-linear 11% 57% 44%)"},
		{"#a64bc8", "color(from #fff srgb-linear 38% 7% 58%)"},
		{"#f5d100", "color(from #fff srgb-linear 91% 64% 0%)"},
		{"#b8acdb", "color(from #fff srgb-linear 48% 41% 71%)"},
		{"#9361f2", "color(from #fff srgb-linear 29% 12% 89%)"},
		{"#ecbda4", "color(from #fff srgb-linear 84% 51% 37%)"},
		{"#dfdd8b", "color(from #fff srgb-linear 74% 72% 26%)"},
		{"#3fcab8", "color(from #fff srgb-linear 5% 59% 48%)"},
		{"#f8a661", "color(from #fff srgb-linear 94% 38% 12%)"},
		{"#bf6f55", "color(from #fff srgb-linear 52% 16% 9%)"},
		{"#818ba8", "color(from #fff srgb-linear 22% 26% 39%)"},

		// keywords
		{"#bad455", "color(from #bad455 srgb-linear r g b)"},
		{"#bad455", "color(from #bad455 srgb-linear r g b / alpha)"},

		// calc(...)
		{"#bad455", "color(from #bad455 srgb-linear calc(r) calc(g) calc(b))"},
		/*{
			"#daa9d56e",
			"color(from #dda0dd srgb-linear calc(r - 0.031) calc(g + 0.052) calc(b * 0.92) / calc(alpha - 0.57))",
		},*/
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
