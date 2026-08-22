package comp

import (
	"fmt"
)

// Linspace returns n points evenly spaced from lo to hi inclusive. It mirrors
// the helper used by the single-component curve so that competitive isotherms
// can be drawn on the same axis.
func Linspace(lo, hi float64, n int) ([]float64, error) {
	if n <= 0 {
		return nil, fmt.Errorf("comp: linspace requires n > 0, got %d", n)
	}
	if hi < lo {
		return nil, fmt.Errorf("comp: linspace requires hi >= lo (got lo=%v hi=%v)", lo, hi)
	}
	xs := make([]float64, n)
	if n == 1 {
		xs[0] = lo
		return xs, nil
	}
	step := (hi - lo) / float64(n-1)
	for i := 0; i < n; i++ {
		xs[i] = lo + float64(i)*step
	}
	return xs, nil
}

// MonotonicA reports whether the coverage of component A is non-decreasing
// across the sampled sweep. A correct competitive Langmuir model is monotone in
// each component's own pressure, so a negative result signals a broken formula.
func MonotonicA(pts []Point) bool {
	for i := 1; i < len(pts); i++ {
		if pts[i].ThetaA < pts[i-1].ThetaA-1e-12 {
			return false
		}
	}
	return true
}

// MaxTotal returns the largest combined coverage seen in a sweep, which is the
// closest approach to the saturation ceiling (still strictly below 1) present
// in the sample.
func MaxTotal(pts []Point) float64 {
	max := 0.0
	for _, p := range pts {
		if p.Total > max {
			max = p.Total
		}
	}
	return max
}

// SumBelowOne confirms the defining competitive invariant: the sum of the two
// coverages never exceeds unity. It returns the largest total found together
// with a boolean that is false the moment any sample breaks the bound.
func SumBelowOne(pts []Point, tol float64) (maxTotal float64, ok bool) {
	ok = true
	for _, p := range pts {
		t := p.ThetaA + p.ThetaB
		if t > maxTotal {
			maxTotal = t
		}
		if t > 1+tol {
			ok = false
		}
	}
	// The analytic form is strictly < 1, so tol is only a numerical safety
	// margin against floating-point rounding at extreme pressures.
	return maxTotal, ok
}
