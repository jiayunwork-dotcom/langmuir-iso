package iso

import "math"

// HalfPressure returns the driving value x_{1/2} at which the coverage equals
// 0.5:
//
//	x_{1/2} = 1 / K.
//
// Doubling K therefore halves the half-coverage pressure, while changing qmax
// leaves it untouched (qmax affects the amount, not the coverage).
func HalfPressure(K float64) float64 {
	if K == 0 {
		return math.Inf(1)
	}
	return 1 / K
}

// HalfCoverageError returns |theta(x_{1/2}) - 0.5| for the supplied K. For a
// correct model it is exactly 0; for a model that has drifted to a different
// saturating form it is positive and can be used as a single scalar diagnostic
// of "is this still a Langmuir isotherm?".
func HalfCoverageError(K float64) float64 {
	return math.Abs(CoverageAtHalf(K) - 0.5)
}
