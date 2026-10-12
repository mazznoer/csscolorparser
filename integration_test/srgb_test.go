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

		// with alpha
		{"#3351fa85", "color(srgb 0.2 0.318 0.979 / 52%)"},
		{"#f08177bf", "color(srgb 0.942 0.504 0.466 / 75%)"},
		{"#fdb5861a", "color(srgb 0.994 0.709 0.525 / 10%)"},
		{"#7d5c58ed", "color(srgb 0.491 0.362 0.344 / 93%)"},
		{"#3168c980", "color(srgb 0.193 0.407 0.79 / 50%)"},
		{"#cc3ff196", "color(srgb 0.8 0.248 0.945 / 59%)"},
		{"#9db3e7f7", "color(srgb 0.615 0.702 0.907 / 97%)"},
		{"#0a414a12", "color(srgb 0.038 0.253 0.29 / 7%)"},
		{"#d6f54eab", "color(srgb 0.838 0.962 0.307 / 67%)"},
		{"#1ad90d59", "color(srgb 0.103 0.852 0.052 / 35%)"},

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

func TestSrgbRelative(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "color(from #bad455 srgb none none none)"},
		{"#bad45500", "color(from #bad455 srgb r g b / none)"},

		// bare number
		{"#9b3af4", "color(from #fff srgb 0.608 0.229 0.958)"},
		{"#4a230f", "color(from #fff srgb 0.29 0.137 0.059)"},
		{"#f2f669", "color(from #fff srgb 0.949 0.966 0.412)"},
		{"#a9389c", "color(from #fff srgb 0.664 0.22 0.612)"},
		{"#edc4b6", "color(from #fff srgb 0.928 0.769 0.712)"},
		{"#db778b", "color(from #fff srgb 0.86 0.468 0.547)"},
		{"#259057", "color(from #fff srgb 0.147 0.566 0.34)"},
		{"#1d5e9f", "color(from #fff srgb 0.113 0.367 0.625)"},
		{"#a936c2", "color(from #fff srgb 0.663 0.212 0.759)"},
		{"#22e2ef", "color(from #fff srgb 0.133 0.885 0.939)"},
		{"#3f7353", "color(from #fff srgb 0.248 0.451 0.325)"},
		{"#c4d753", "color(from #fff srgb 0.768 0.844 0.326)"},
		{"#d5a488", "color(from #fff srgb 0.835 0.645 0.532)"},
		{"#d10f04", "color(from #fff srgb 0.819 0.06 0.014)"},
		{"#5445a5", "color(from #fff srgb 0.331 0.272 0.647)"},

		// percentage
		{"#696e3d", "color(from #fff srgb 41% 43% 24%)"},
		{"#572ee0", "color(from #fff srgb 34% 18% 88%)"},
		{"#c2d170", "color(from #fff srgb 76% 82% 44%)"},
		{"#ed29cc", "color(from #fff srgb 93% 16% 80%)"},
		{"#03d9e8", "color(from #fff srgb 1% 85% 91%)"},
		{"#4a4273", "color(from #fff srgb 29% 26% 45%)"},
		{"#bf4773", "color(from #fff srgb 75% 28% 45%)"},
		{"#3be3c4", "color(from #fff srgb 23% 89% 77%)"},
		{"#572121", "color(from #fff srgb 34% 13% 13%)"},
		{"#75bfe3", "color(from #fff srgb 46% 75% 89%)"},
		{"#803b69", "color(from #fff srgb 50% 23% 41%)"},
		{"#4d6bfa", "color(from #fff srgb 30% 42% 98%)"},
		{"#d4d40d", "color(from #fff srgb 83% 83% 5%)"},
		{"#bfa330", "color(from #fff srgb 75% 64% 19%)"},
		{"#2eed54", "color(from #fff srgb 18% 93% 33%)"},

		// keywords
		{"#bad455", "color(from #bad455 srgb r g b)"},
		{"#bad455", "color(from #bad455 srgb r g b / alpha)"},

		// calc(...)
		{"#d455ba", "color(from #bad455 srgb calc(g) calc(b) calc(r))"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
