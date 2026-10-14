package csscolorparser

import "math"

// D50 reference white (CIE xy = 0.3457, 0.3585), as used by CSS Color 4.
const (
	d50X = 0.3457 / 0.3585
	d50Y = 1.0
	d50Z = (1.0 - 0.3457 - 0.3585) / 0.3585
)

// CIE constants for Lab.
const (
	epsilon = 216.0 / 24389.0 // (6/29)^3
	kappa   = 24389.0 / 27.0  // (29/3)^3
)

// linearRGBToXYZD50 converts linear-light sRGB (components in 0..1, no gamma)
// to CIE XYZ relative to the D50 white point.
//
// The matrix is sRGB primaries (D65) chromatically adapted to D50 with the
// Bradford transform.
func linearRGBToXYZD50(r, g, b float64) (x, y, z float64) {
	x = 0.4360747*r + 0.3850649*g + 0.1430804*b
	y = 0.2225045*r + 0.7168786*g + 0.0606169*b
	z = 0.0139322*r + 0.0971045*g + 0.7141733*b
	return
}

// xyzD50ToLinearRGB converts CIE XYZ (D50) to linear-light sRGB.
// It is the inverse of linearRGBToXYZD50. Results are not clamped, so
// out-of-gamut colors yield components below 0 or above 1.
func xyzD50ToLinearRGB(x, y, z float64) (r, g, b float64) {
	r = 3.1338564*x - 1.6168668*y - 0.4906148*z
	g = -0.9787686*x + 1.9161416*y + 0.0334541*z
	b = 0.0719452*x - 0.2289913*y + 1.4052427*z
	return
}

// xyzD50ToLab converts CIE XYZ (D50) to CIE Lab (D50).
// L is in 0..100; a and b are roughly in -125..125.
func xyzD50ToLab(x, y, z float64) (l, a, b float64) {
	fx := labF(x / d50X)
	fy := labF(y / d50Y)
	fz := labF(z / d50Z)

	l = 116*fy - 16
	a = 500 * (fx - fy)
	b = 200 * (fy - fz)
	return
}

func labF(t float64) float64 {
	if t > epsilon {
		return math.Cbrt(t)
	}
	return (kappa*t + 16) / 116
}

// labToXYZD50 converts CIE Lab (D50) to CIE XYZ relative to the D50 white point.
func labToXYZD50(l, a, b float64) (x, y, z float64) {
	fy := (l + 16) / 116
	fx := a/500 + fy
	fz := fy - b/200

	fx3 := fx * fx * fx
	fz3 := fz * fz * fz

	if fx3 > epsilon {
		x = fx3
	} else {
		x = (116*fx - 16) / kappa
	}
	if l > kappa*epsilon {
		y = fy * fy * fy
	} else {
		y = l / kappa
	}
	if fz3 > epsilon {
		z = fz3
	} else {
		z = (116*fz - 16) / kappa
	}

	return x * d50X, y * d50Y, z * d50Z
}

// linearRGBToLab is a convenience wrapper chaining linear RGB -> XYZ D50 -> Lab.
func linearRGBToLab(r, g, b float64) (L, A, B float64) {
	x, y, z := linearRGBToXYZD50(r, g, b)
	return xyzD50ToLab(x, y, z)
}

// labToLinearRGB converts CIE Lab (D50) to linear-light sRGB by chaining
// Lab -> XYZ D50 -> linear RGB. Results are not clamped.
func labToLinearRGB(l, a, b float64) (R, G, B float64) {
	x, y, z := labToXYZD50(l, a, b)
	return xyzD50ToLinearRGB(x, y, z)
}
