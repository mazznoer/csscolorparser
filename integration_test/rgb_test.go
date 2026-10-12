package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestRgbRelative(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "rgb(from #fff none none none)"},
		{"#bad45500", "rgb(from #bad455 r g b / none)"},

		// bare number
		{"#df5d70", "rgb(from #fff 223 93 112)"},
		{"#8d4b7b", "rgb(from #fff 141 75 123)"},
		{"#b9615f", "rgb(from #fff 185 97 95)"},
		{"#afc7bf", "rgb(from #fff 175 199 191)"},
		{"#9b6bcf", "rgb(from #fff 155 107 207)"},
		{"#1dc57b", "rgb(from #fff 29 197 123)"},
		{"#02349a", "rgb(from #fff 2 52 154)"},
		{"#a83caf", "rgb(from #fff 168 60 175)"},
		{"#9c77f0", "rgb(from #fff 156 119 240)"},
		{"#9e533a", "rgb(from #fff 158 83 58)"},
		{"#97eb78", "rgb(from #fff 151 235 120)"},
		{"#c8a9da", "rgb(from #fff 200 169 218)"},
		{"#6b6021", "rgb(from #fff 107 96 33)"},
		{"#9accfc", "rgb(from #fff 154 204 252)"},
		{"#a2dd3f", "rgb(from #fff 162 221 63)"},
		{"#34a783", "rgb(from #fff 52 167 131)"},

		// percentage
		{"#e85447", "rgb(from #fff 91% 33% 28%)"},
		{"#63b8fa", "rgb(from #fff 39% 72% 98%)"},
		{"#d6e09e", "rgb(from #fff 84% 88% 62%)"},
		{"#d4d445", "rgb(from #fff 83% 83% 27%)"},
		{"#0fa340", "rgb(from #fff 6% 64% 25%)"},
		{"#5c7d42", "rgb(from #fff 36% 49% 26%)"},
		{"#00cf69", "rgb(from #fff 0% 81% 41%)"},
		{"#17e3d4", "rgb(from #fff 9% 89% 83%)"},
		{"#5e6682", "rgb(from #fff 37% 40% 51%)"},
		{"#ccad69", "rgb(from #fff 80% 68% 41%)"},
		{"#c75221", "rgb(from #fff 78% 32% 13%)"},
		{"#e3f01a", "rgb(from #fff 89% 94% 10%)"},
		{"#52385c", "rgb(from #fff 32% 22% 36%)"},
		{"#9cad7d", "rgb(from #fff 61% 68% 49%)"},
		{"#c229ba", "rgb(from #fff 76% 16% 73%)"},
		{"#0fe8d9", "rgb(from #fff 6% 91% 85%)"},

		// keywords
		{"#bad455", "rgb(from #bad455 r g b)"},
		{"#55bad4", "rgb(from #bad455 b r g / alpha)"},

		// calc(...)
		{"#bad455", "rgb(from #bad455 calc(r) calc(g) calc(b))"},
		{"#a18ac5", "rgb(from #bad455 calc((r + g + b) / 3) calc(12 * 9 + 30) calc(0.197 * 1000))"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
