package csscolorparser

import (
	"math"
	"strconv"
	"strings"
)

func modulo(x, y float64) float64 {
	return math.Mod(math.Mod(x, y)+y, y)
}

func clamp(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func parseFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f, err == nil
}

// Returns (result, ok?, percentage?)
func parsePercentOrFloat(s string) (float64, bool, bool) {
	if strings.HasSuffix(s, "%") {
		f, ok := parseFloat(s[:len(s)-1])
		if ok {
			return f / 100, true, true
		}
		return 0, false, true
	}
	f, ok := parseFloat(s)
	return f, ok, false
}

// Result angle in degrees (not normalized)
func parseAngle(s string) (float64, bool) {
	if strings.HasSuffix(s, "deg") {
		return parseFloat(s[:len(s)-3])
	}
	if strings.HasSuffix(s, "grad") {
		f, ok := parseFloat(s[:len(s)-4])
		if ok {
			return f / 400 * 360, true
		}
		return 0, false
	}
	if strings.HasSuffix(s, "rad") {
		f, ok := parseFloat(s[:len(s)-3])
		if ok {
			return f / math.Pi * 180, true
		}
		return 0, false
	}
	if strings.HasSuffix(s, "turn") {
		f, ok := parseFloat(s[:len(s)-4])
		if ok {
			return f * 360, true
		}
		return 0, false
	}
	return parseFloat(s)
}

func normalizeAngle(t float64) float64 {
	t = math.Mod(t, 360)
	if t < 0 {
		t += 360
	}
	return t
}

// Map t which is in range [a, b] to range [c, d]
func remap(t, a, b, c, d float64) float64 {
	return (t-a)*((d-c)/(b-a)) + c
}
