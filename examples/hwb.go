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
	w := 400
	h := 400
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	fw := float64(w)
	fh := float64(h)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			s := fmt.Sprintf("hwb(120 %v%% %v%%)", remap(float64(y), 0, fh, 0, 100), remap(float64(x), 0, fw, 0, 100))
			col, _ := csscolorparser.Parse(s)
			img.Set(x, y, col)
		}
	}

	file, err := os.Create("hwb-3.png")

	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	png.Encode(file, img)
}

func remap(value, a, b, c, d float64) float64 {
	return (value-a)*((d-c)/(b-a)) + c
}
