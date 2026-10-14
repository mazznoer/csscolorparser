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

func TestOklabRelative(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "oklab(from #bad455 none none none)"},
		{"#bad45500", "oklab(from #bad455 l a b / none)"},

		// bare number
		{"#32340c", "oklab(from #fff 0.314 -0.022 0.054)"},
		{"#00407e", "oklab(from #fff 0.297 -0.316 -0.149)"},
		{"#003918", "oklab(from #fff 0.279 -0.118 0.032)"},
		{"#111f00", "oklab(from #fff 0.213 -0.13 0.222)"},
		{"#003adc", "oklab(from #fff 0.383 -0.287 -0.278)"},
		{"#00caff", "oklab(from #fff 0.783 -0.195 -0.261)"},
		{"#ff3100", "oklab(from #fff 0.946 0.255 0.368)"},
		{"#9300bd", "oklab(from #fff 0.397 0.294 -0.243)"},
		{"#ffaac2", "oklab(from #fff 0.917 0.184 0.039)"},
		{"#006400", "oklab(from #fff 0.402 -0.208 0.14)"},
		{"#2d007b", "oklab(from #fff 0.292 0.047 -0.166)"},
		{"#aa0000", "oklab(from #fff 0.275 0.32 0.344)"},
		{"#edbc00", "oklab(from #fff 0.812 -0.017 0.208)"},
		{"#00b8ff", "oklab(from #fff 0.709 -0.244 -0.257)"},
		{"#8a0000", "oklab(from #fff 0.37 0.138 0.15)"},

		// percentage
		{"#26007f", "oklab(from #fff 10% -52% -74%)"},
		{"#001100", "oklab(from #fff 11% -55% 20%)"},
		{"#ffb100", "oklab(from #fff 96% 36% 66%)"},
		{"#d5d2ff", "oklab(from #fff 88% 5% -15%)"},
		{"#00a3ff", "oklab(from #fff 80% -18% -91%)"},
		{"#390098", "oklab(from #fff 13% 41% -88%)"},
		//{"#ffa3da", "oklab(from #fff 93% 53% 3%)"},
		{"#00c000", "oklab(from #fff 66% -74% 53%)"},
		{"#dd0000", "oklab(from #fff 51% 48% 64%)"},
		{"#00e6ff", "oklab(from #fff 86% -87% -97%)"},
		{"#86bd00", "oklab(from #fff 72% -58% 97%)"},
		{"#00a500", "oklab(from #fff 57% -75% 40%)"},
		{"#0b0030", "oklab(from #fff 2% -14% -45%)"},
		{"#0061ff", "oklab(from #fff 50% -61% -74%)"},
		{"#4f0000", "oklab(from #fff 12% 74% 37%)"},

		// keywords
		{"#bad455", "oklab(from #bad455 l a b)"},
		{"#bad455", "oklab(from #bad455 l a b / alpha)"},

		// calc(...)
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
