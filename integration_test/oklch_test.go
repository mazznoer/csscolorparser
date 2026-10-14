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

func TestOklchRelative(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "oklch(from #bad455 none none none)"},
		{"#bad45500", "oklch(from #bad455 l c h / none)"},

		// bare number
		{"#ff80be", "oklch(from #fff 0.818 0.208 0)"},
		{"#610000", "oklch(from #fff 0.252 0.193 25)"},
		{"#ff0000", "oklch(from #fff 0.75 0.38 51)"},
		//{"#ffb532", "oklch(from #fff 0.822 0.159 76)"},
		{"#806400", "oklch(from #fff 0.509 0.209 102)"},
		{"#beff49", "oklch(from #fff 0.931 0.217 128)"},
		{"#13cc6f", "oklch(from #fff 0.743 0.188 153)"},
		{"#00432c", "oklch(from #fff 0.286 0.207 179)"},
		{"#009eac", "oklch(from #fff 0.634 0.117 205)"},
		{"#00dcff", "oklch(from #fff 0.824 0.262 230)"},
		{"#220097", "oklch(from #fff 0.182 0.305 256)"},
		{"#c1c4ff", "oklch(from #fff 0.847 0.111 282)"},
		{"#a848ef", "oklch(from #fff 0.607 0.24 307)"},
		{"#ff5fff", "oklch(from #fff 0.778 0.271 333)"},
		{"#ffd9f8", "oklch(from #fff 0.985 0.116 359)"},

		// percentage
		{"#004439", "oklch(from #fff 34% 18% 179)"},
		{"#200059", "oklch(from #fff 9% 69% 299)"},
		{"#004bb1", "oklch(from #fff 43% 47% 251)"},
		{"#ffb6ca", "oklch(from #fff 88% 30% 7)"},
		{"#e7bb4c", "oklch(from #fff 81% 34% 87)"},
		{"#7e7572", "oklch(from #fff 57% 3% 43)"},
		{"#822b31", "oklch(from #fff 42% 30% 20)"},
		{"#920000", "oklch(from #fff 39% 49% 33)"},
		{"#5e33d2", "oklch(from #fff 48% 56% 287)"},
		{"#2f556b", "oklch(from #fff 43% 14% 236)"},
		{"#0095d8", "oklch(from #fff 56% 79% 210)"},
		{"#1d0000", "oklch(from #fff 12% 22% 17)"},
		{"#d4e200", "oklch(from #fff 87% 56% 114)"},
		{"#00ffab", "oklch(from #fff 88% 61% 164)"},
		{"#33424a", "oklch(from #fff 37% 6% 231)"},

		// keywords
		{"#bad455", "oklch(from #bad455 l c h)"},
		{"#bad455", "oklch(from #bad455 l c h / alpha)"},

		// calc(...)
		{"#bad455", "oklch(from #bad455 calc(l) calc(c) calc(h))"},
		{"#7aa800", "oklch(from #bad455 calc(0.5 + 0.17) calc(1 * 0.2) calc(10 * 12.5))"},
		{"#ff7e00", "oklch(from #bad455 calc(l - 0.1) calc(c + 0.05) calc(h / 2))"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
