package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestHwbAbsolute(t *testing.T) {
	testData := [][2]string{
		// none
		{"#ff0000", "hwb(none none none)"},
		{"#00b7ff00", "hwb(197 0 0 / none)"},

		// bare number
		{"#420000", "hwb(0 0 74)"},
		{"#692800", "hwb(23 0 59)"},
		{"#ffdd61", "hwb(47 38 0)"},
		{"#89a800", "hwb(71 0 34)"},
		{"#9bff54", "hwb(95 33 0)"},
		{"#03a300", "hwb(119 0 36)"},
		{"#00b846", "hwb(143 0 28)"},
		{"#003d30", "hwb(167 0 76)"},
		{"#baf2ff", "hwb(191 73 0)"},
		{"#78b0ff", "hwb(215 47 0)"},
		{"#babbff", "hwb(239 73 0)"},
		{"#4300b0", "hwb(263 0 31)"},
		{"#5e0078", "hwb(287 0 53)"},
		{"#ff1ad5", "hwb(311 10 0)"},
		{"#260010", "hwb(335 0 85)"},

		// angle suffix
		{"#eee3ff", "hwb(263deg 89 0)"},
		{"#ff8cf2", "hwb(341grad 55 0)"},
		{"#8a4a00", "hwb(0.090turn 0 46)"},
		{"#787bff", "hwb(265grad 47 0)"},
		{"#ddffab", "hwb(1.472rad 67 0)"},
		{"#001421", "hwb(227grad 0 87)"},
		{"#ca00de", "hwb(5.144rad 0 13)"},
		{"#92ff8c", "hwb(0.325turn 55 0)"},
		{"#fffd33", "hwb(0.165turn 20 0)"},
		{"#ff858e", "hwb(6.200rad 52 0)"},
		{"#b0ecff", "hwb(3.399rad 69 0)"},
		{"#13001c", "hwb(0.782turn 0 89)"},
		{"#001f1e", "hwb(3.109rad 0 88)"},
		{"#f8ffd6", "hwb(0.194turn 84 0)"},
		{"#00a368", "hwb(0.440turn 0 36)"},
		{"#302800", "hwb(0.136turn 0 81)"},
		{"#5c009c", "hwb(0.765turn 0 39)"},
		{"#a3ffb7", "hwb(148grad 64 0)"},
		{"#4f85ff", "hwb(0.616turn 31 0)"},
		{"#0c3600", "hwb(1.851rad 0 79)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
