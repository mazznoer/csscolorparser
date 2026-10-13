package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestOklabAbsolute(t *testing.T) {
	testData := [][2]string{
		// none
		{"#00000000", "oklab(none none none / none)"},

		// bare number
		{"#4f00ff", "oklab(0.476 0.043 -0.334)"},
		{"#e0ff00", "oklab(1 -0.163 0.264)"},
		{"#e700ff", "oklab(0.661 0.274 -0.246)"},
		{"#778a91", "oklab(0.621 -0.019 -0.016)"},
		{"#b7c2ff", "oklab(0.844 0.02 -0.138)"},
		{"#003e49", "oklab(0.329 -0.06 -0.036)"},
		{"#000070", "oklab(0.147 -0.338 -0.219)"},
		{"#ee0000", "oklab(0.387 0.368 0.396)"},
		{"#005800", "oklab(0.36 -0.343 0.389)"},
		{"#183000", "oklab(0.271 -0.215 0.366)"},
		{"#f40000", "oklab(0.527 0.215 0.311)"},
		{"#8c79cc", "oklab(0.628 0.048 -0.114)"},
		{"#00f9bc", "oklab(0.81 -0.306 0.018)"},
		{"#dd9900", "oklab(0.724 0.002 0.223)"},
		{"#6600c6", "oklab(0.425 0.119 -0.232)"},

		// percentage
		{"#006000", "oklab(41% -47% 53%)"},
		{"#ffb800", "oklab(87% 13% 60%)"},
		{"#ffa2fb", "oklab(88% 44% -15%)"},
		{"#ff00b8", "oklab(66% 86% -19%)"},
		{"#00e6b1", "oklab(73% -94% 0%)"},
		{"#005179", "oklab(39% -25% -23%)"},
		{"#4200ee", "oklab(29% 4% -99%)"},
		{"#956200", "oklab(53% -6% 58%)"},
		{"#ffd2ff", "oklab(98% 28% -45%)"},
		{"#001fff", "oklab(51% -53% -100%)"},
		{"#150055", "oklab(9% -24% -54%)"},
		{"#009af3", "oklab(61% -49% -45%)"},
		{"#001300", "oklab(14% -44% 47%)"},
		{"#57a225", "oklab(64% -31% 30%)"},
		{"#00ffff", "oklab(90% -78% -37%)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
