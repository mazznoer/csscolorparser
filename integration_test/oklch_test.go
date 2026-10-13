package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestOklchAbsolute(t *testing.T) {
	testData := [][2]string{
		// none
		{"#00000000", "oklch(none none none / none)"},

		// bare number
		{"#ff00a7", "oklch(0.759 0.397 0)"},
		{"#ff0000", "oklch(0.73 0.372 25)"},
		{"#741a00", "oklch(0.366 0.142 51)"},
		{"#8b0000", "oklch(0.346 0.331 76)"},
		{"#140500", "oklch(0.112 0.247 102)"},
		{"#70cc00", "oklch(0.749 0.318 128)"},
		{"#009200", "oklch(0.511 0.348 153)"},
		{"#00d9b2", "oklch(0.754 0.219 179)"},
		{"#86f6ff", "oklch(0.912 0.104 205)"},
		{"#00ffff", "oklch(0.96 0.231 230)"},
		{"#00c8ff", "oklch(0.87 0.337 256)"},
		{"#8a7cff", "oklch(0.686 0.233 282)"},
		{"#4a0096", "oklch(0.222 0.354 307)"},
		{"#e476d4", "oklch(0.721 0.177 333)"},
		{"#ff00e0", "oklch(0.924 0.385 359)"},

		// percentage + angle prefix
		{"#e7e6d4", "oklch(92% 6% 116grad)"},
		{"#424835", "oklch(39% 8% 0.335turn)"},
		{"#007278", "oklch(50% 22% 224grad)"},
		{"#5fdbff", "oklch(84% 32% 252grad)"},
		{"#4900e2", "oklch(28% 98% 312grad)"},
		{"#ffc0f8", "oklch(96% 42% 355deg)"},
		{"#ffb14e", "oklch(86% 43% 58deg)"},
		{"#ff00a1", "oklch(87% 91% 10grad)"},
		{"#00bdff", "oklch(70% 88% 220deg)"},
		{"#001f00", "oklch(11% 98% 1.969rad)"},
		{"#a8a8a2", "oklch(73% 2% 107deg)"},
		{"#6e71ff", "oklch(65% 63% 276deg)"},
		{"#00d646", "oklch(72% 77% 157deg)"},
		{"#ff0088", "oklch(58% 96% 0.984turn)"},
		{"#00daff", "oklch(74% 93% 3.602rad)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
