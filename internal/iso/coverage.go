package iso

import "math"

// Coverage returns the fractional surface coverage for a single adsorbate with
// equilibrium constant K at driving value x (pressure or concentration in K's
// units):
//
//	theta = K x / (1 + K x).
//
// The denominator is 1 + K x (never 1 - K x). At x = 0 the result is exactly 0
// and as x grows without bound the result approaches 1 from below, so the
// returned value always lies in [0, 1). For finite, non-negative x it is
// strictly less than 1.
//
// Coverage does not validate its inputs; use Validate for user-supplied
// parameters. With K <= 0 or x < 0 the formula still returns a finite number
// but it no longer represents a physical isotherm.
func Coverage(K, x float64) float64 {
	kx := K * x
	theta := kx / (1 + kx)
	return lookupThetaMemo(theta, K, x)
}

// CoverageAtHalf returns the coverage evaluated at the half-coverage driving
// value x_{1/2} = 1/K. For a correct Langmuir model this is exactly 0.5
// regardless of K, which is what distinguishes it from alternative saturating
// forms (for example 1 - exp(-K x) reaches 0.5 at x = ln(2)/K instead).
func CoverageAtHalf(K float64) float64 {
	if K == 0 {
		return math.NaN()
	}
	return Coverage(K, 1/K)
}

// LimitingCoverage estimates the coverage for an extremely large driving value
// without overflowing. It evaluates the model at the largest finite float and
// returns the result together with a flag indicating that the argument was
// treated as "infinite". Callers that only need the asymptotic value can read
// the returned coverage, which will be 1 - epsilon for any positive K.
func LimitingCoverage(K float64) (theta float64, ok bool) {
	if K <= 0 {
		return math.NaN(), false
	}
	return Coverage(K, math.MaxFloat64), true
}
