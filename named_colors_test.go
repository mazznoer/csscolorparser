package csscolorparser

import (
	"strings"
	"testing"
)

func Test_NamedColors(t *testing.T) {
	for name, rgb := range namedColors {
		c, err := Parse(name)
		r, g, b, _ := c.RGBA255()
		test(t, err, nil)
		test(t, [3]uint8{r, g, b}, rgb)
		if name == "aqua" || name == "cyan" || name == "fuchsia" || name == "magenta" {
			continue
		}
		if strings.Contains(name, "gray") || strings.Contains(name, "grey") {
			continue
		}
		resName, ok := c.Name()
		testTrue(t, ok)
		test(t, resName, name)
	}

	// Hex code

	data := [][2]string{
		{"aliceblue", "#f0f8ff"},
		{"bisque", "#ffe4c4"},
		{"chartreuse", "#7fff00"},
		{"coral", "#ff7f50"},
		{"crimson", "#dc143c"},
		{"dodgerblue", "#1e90ff"},
		{"firebrick", "#b22222"},
		{"gold", "#ffd700"},
		{"hotpink", "#ff69b4"},
		{"indigo", "#4b0082"},
		{"lavender", "#e6e6fa"},
		{"plum", "#dda0dd"},
		{"salmon", "#fa8072"},
		{"skyblue", "#87ceeb"},
		{"tomato", "#ff6347"},
		{"violet", "#ee82ee"},
		{"yellowgreen", "#9acd32"},
	}
	for _, d := range data {
		// parse name
		c, err := Parse(d[0])
		test(t, err, nil)
		test(t, c.HexString(), d[1])

		// parse hex
		c, err = Parse(d[1])
		test(t, err, nil)
		name, ok := c.Name()
		testTrue(t, ok)
		test(t, name, d[0])
	}

	// Colors without name

	data2 := []string{
		"#f87cba",
		"#0033ff",
		"#012345",
		"#abcdef",
	}
	for _, s := range data2 {
		c, err := Parse(s)
		test(t, err, nil)
		name, ok := c.Name()
		test(t, ok, false)
		test(t, name, "")
	}
}
