// +build ignore

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/mazznoer/csscolorparser"
)

func main() {
	w := 800
	h := 100
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	fw := float64(w)
	fh := float64(h)

	for x := 0; x < w; x++ {
		//s := fmt.Sprintf("hsl(%v,100%%,50%%)", remap(float64(x), 0, fw, 0, 360))
		//col, _ := csscolorparser.Parse(s)

		for y := 0; y < h; y++ {
			s := fmt.Sprintf("hsl(%v,100%%,%v%%)", remap(float64(x), 0, fw, 0, 360), remap(float64(y), 0, fh, 100, 0))
			col, _ := csscolorparser.Parse(s)
			img.Set(x, y, col)
		}
	}

	file, err := os.Create("rainbow-2.png")

	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	png.Encode(file, img)
}

func remap(value, a, b, c, d float64) float64 {
	return (value-a)*((d-c)/(b-a)) + c
}
