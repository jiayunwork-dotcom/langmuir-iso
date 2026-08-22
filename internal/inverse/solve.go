package inverse

import (
	"fmt"
	"math"
)

// PressureFromTheta returns the driving value x that yields coverage theta for
// equilibrium constant K:
//
//	x = theta / (K (1 - theta)).
//
// The inverse is undefined at theta = 1 (full coverage is only reached as x ->
// +infinity) and at theta >= 1 or theta < 0 the coverage is unphysical, so the
// function returns an error in those cases instead of a spurious number.
//
// The implementation guards against divide-by-zero explicitly and against
// non-finite inputs, so it never panics on the boundary values.
func PressureFromTheta(K, theta float64) (float64, error) {
	if !isFinite(K) || !isFinite(theta) {
		return 0, fmt.Errorf("inverse: k and theta must be finite (got k=%v theta=%v)", K, theta)
	}
	if K <= 0 {
		return 0, fmt.Errorf("inverse: k must be > 0 (got %v)", K)
	}
	if theta < 0 {
		return 0, fmt.Errorf("inverse: theta must be >= 0 (got %v)", theta)
	}
	if theta >= 1 {
		return 0, fmt.Errorf("inverse: theta=1 has no finite driving value (got theta=%v)", theta)
	}
	denom := K * (1 - theta)
	if denom == 0 {
		return 0, fmt.Errorf("inverse: degenerate denominator for theta=%v", theta)
	}
	return theta / denom, nil
}

// Solve runs PressureFromTheta for the structured problem and wraps the result
// in a Solution for callers that prefer to carry the whole query around.
func Solve(p Problem) (Solution, error) {
	x, err := PressureFromTheta(p.K, p.Theta)
	if err != nil {
		return Solution{}, err
	}
	return Solution{K: p.K, Theta: p.Theta, X: x}, nil
}

// Check recovers the coverage from a solved driving value and compares it to the
// requested theta. A correct solver round-trips: Coverage(K, x) == theta to
// within tol. This is the public invariant the tests assert.
func (s Solution) Check(tol float64) (recovered float64, ok bool) {
	kx := s.K * s.X
	recovered = kx / (1 + kx)
	return recovered, math.Abs(recovered-s.Theta) <= tol
}

func isFinite(v float64) bool {
	return v == v && v != math.Inf(1) && v != math.Inf(-1)
}
