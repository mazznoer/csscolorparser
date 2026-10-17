package integration_test

import (
	"iter"
	"testing"

	"github.com/mazznoer/csscolorparser"
)

func TestParseColors(t *testing.T) {
	s := "red, #00ff00, rgb(0, 0, 255), #x00, #bad455"
	var colors []string

	for c, err := range csscolorparser.ParseColors(s) {
		if err != nil {
			continue
		}
		colors = append(colors, c.HexString())
	}
	test(t, len(colors), 4)

	// ---

	next, stop := iter.Pull2(csscolorparser.ParseColors(s))
	defer stop()

	c, err, ok := next()
	testTrue(t, ok)
	testTrue(t, err == nil)
	test(t, c.HexString(), "#ff0000")

	c, err, ok = next()
	testTrue(t, ok)
	testTrue(t, err == nil)
	test(t, c.HexString(), "#00ff00")

	c, err, ok = next()
	testTrue(t, ok)
	testTrue(t, err == nil)
	test(t, c.HexString(), "#0000ff")

	c, err, ok = next()
	testTrue(t, ok)
	testTrue(t, err != nil)

	c, err, ok = next()
	testTrue(t, ok)
	testTrue(t, err == nil)
	test(t, c.HexString(), "#bad455")

	c, err, ok = next()
	testTrue(t, !ok)
}
