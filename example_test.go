package csscolorparser_test

import (
	"fmt"

	"github.com/mazznoer/csscolorparser"
)

func Example_namedColor() {
	c, err := csscolorparser.Parse("gold")

	if err != nil {
		panic(err)
	}

	fmt.Printf("R:%.3f, G:%.3f, B:%.3f, A:%.3f\n", c.R, c.G, c.B, c.A)
	fmt.Println(c.RGBA255())
	fmt.Println(c.HexString())
	fmt.Println(c.RGBString())
	// Output:
	// R:1.000, G:0.843, B:0.000, A:1.000
	// 255 215 0 255
	// #ffd700
	// rgb(255,215,0)
}

func Example_rgbColor() {
	c, err := csscolorparser.Parse("rgba(100%, 0%, 0%, 0.5)")

	if err != nil {
		panic(err)
	}

	fmt.Printf("R:%.3f, G:%.3f, B:%.3f, A:%.3f\n", c.R, c.G, c.B, c.A)
	fmt.Println(c.RGBA255())
	fmt.Println(c.HexString())
	fmt.Println(c.RGBString())
	// Output:
	// R:1.000, G:0.000, B:0.000, A:0.500
	// 255 0 0 128
	// #ff000080
	// rgba(255,0,0,0.5)
}
