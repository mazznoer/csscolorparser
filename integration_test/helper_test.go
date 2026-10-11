package integration_test

import (
	"image/color"
	"testing"

	"github.com/mazznoer/csscolorparser"
)

// --- Helper functions

func test(t *testing.T, a, b interface{}) {
	if a != b {
		t.Helper()
		t.Errorf("left: %v, right: %v", a, b)
	}
}

func testTrue(t *testing.T, b bool) {
	if !b {
		t.Helper()
		t.Errorf("it false")
	}
}

func testColor(t *testing.T, a, b csscolorparser.Color) {
	x := arr8(a.RGBA255())
	y := arr8(b.RGBA255())
	if x == y {
		return
	}
	t.Helper()
	t.Errorf("left: %v, right: %v", x, y)
}

func testGoColor(t *testing.T, a, b color.Color) {
	x := arr32(a.RGBA())
	y := arr32(b.RGBA())
	if x == y {
		return
	}
	t.Helper()
	t.Errorf("left: %v, right: %v", x, y)
}

func arr8(r, g, b, a uint8) [4]uint8 {
	return [4]uint8{r, g, b, a}
}

func arr32(r, g, b, a uint32) [4]uint32 {
	return [4]uint32{r, g, b, a}
}
