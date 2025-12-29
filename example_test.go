package csscolorparser_test

import (
	"fmt"

	"github.com/mazznoer/csscolorparser"
)

func Example() {
	c, err := csscolorparser.Parse("rgb(100% 0% 0% / 50%)")

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
	// rgb(255 0 0 / 50%)
}

func ExampleColor_Name() {
	c, err := csscolorparser.Parse("#ffd700")

	if err != nil {
		panic(err)
	}

	fmt.Println(c.RGBString())
	name, _ := c.Name()
	fmt.Println(name)
	// Output:
	// rgb(255 215 0)
	// gold
}
