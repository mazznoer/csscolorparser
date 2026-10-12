package integration_test

import (
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestHslAbsolute(t *testing.T) {
	testData := [][2]string{
		// none
		{"#000000", "hsl(none none none)"},
		{"#00000000", "hsl(197 0 0 / none)"},

		// bare number
		{"#421f1f", "hsl(0 37 19)"},
		{"#88481b", "hsl(25 67 32)"},
		{"#b49f27", "hsl(51 64 43)"},
		{"#2f3a13", "hsl(76 51 15)"},
		{"#2a481e", "hsl(102 42 20)"},
		{"#3cf655", "hsl(128 91 60)"},
		{"#8ce3bc", "hsl(153 61 72)"},
		{"#388584", "hsl(179 41 37)"},
		{"#2f5a79", "hsl(205 44 33)"},
		{"#17256e", "hsl(230 66 26)"},
		{"#815bec", "hsl(256 79 64)"},
		{"#dfa5f8", "hsl(282 86 81)"},
		{"#6a2963", "hsl(307 44 29)"},
		{"#d8c0cb", "hsl(333 24 80)"},
		{"#fdfcfc", "hsl(359 11 99)"},

		// angle suffix
		{"#1c415a", "hsl(0.565turn 53 23)"},
		{"#dceed8", "hsl(120grad 40 89)"},
		{"#af088d", "hsl(347grad 91 36)"},
		{"#ebfed2", "hsl(0.239turn 95 91)"},
		{"#ab6d93", "hsl(323deg 27 55)"},
		{"#433743", "hsl(0.842turn 10 24)"},
		{"#cf2639", "hsl(0.981turn 69 48)"},
		{"#c7a957", "hsl(44deg 50 56)"},
		{"#3aa305", "hsl(100deg 94 33)"},
		{"#45565f", "hsl(0.557turn 16 32)"},
		{"#76919e", "hsl(221grad 17 54)"},
		{"#5f3321", "hsl(20grad 49 25)"},
		{"#c9cfab", "hsl(1.197rad 27 74)"},
		{"#91baf8", "hsl(216deg 88 77)"},
		{"#fce9e9", "hsl(0.003turn 76 95)"},
		{"#43143c", "hsl(343grad 54 17)"},
		{"#a4d0ab", "hsl(143grad 32 73)"},
		{"#0e565d", "hsl(185deg 73 21)"},
		{"#e46caf", "hsl(363grad 69 66)"},
		{"#a0caa7", "hsl(2.281rad 29 71)"},
		{"#19578f", "hsl(3.640rad 70 33)"},
		{"#ebd57f", "hsl(0.838rad 73 71)"},
		{"#cc75b9", "hsl(0.870turn 46 63)"},
		{"#c8a4d6", "hsl(4.937rad 38 74)"},
		{"#532228", "hsl(0.980turn 42 23)"},
		{"#40b55f", "hsl(2.371rad 48 48)"},
		{"#def575", "hsl(1.236rad 86 71)"},
		{"#15291d", "hsl(2.539rad 33 12)"},
	}

	for _, pair := range testData {
		hex, s := pair[0], pair[1]
		c, err := csscolorparser.Parse(s)
		test(t, err, nil)
		test(t, hex, c.HexString())
	}
}
