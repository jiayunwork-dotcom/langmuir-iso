package iso

import "math"

// HenrySlope returns the initial slope of the coverage curve in the limit of
// vanishing driving value. For the Langmuir isotherm
//
//	theta(x) = K x / (1 + K x) = K x - (K x)^2 + ...
//
// so theta(x)/x -> K as x -> 0. The Henry constant is therefore exactly K.
func HenrySlope(K float64) float64 {
	return K
}

// HenryEstimate returns the first-order Henry approximation of the coverage at a
// small driving value x:
//
//	theta ~= K x.
//
// It is provided so callers can compare the exact model against its low-pressure
// asymptote directly. It is only accurate when K x << 1.
func HenryEstimate(K, x float64) float64 {
	return K * x
}

// HenryRelativeError returns |theta(x) - K x| / (K x) for a small x. It is the
// relative error made by trusting the Henry law alone and shrinks linearly with
// x. Callers use it to decide whether the low-pressure approximation is good
// enough for a given sample.
func HenryRelativeError(K, x float64) float64 {
	if K == 0 || x == 0 {
		return 0
	}
	exact := Coverage(K, x)
	approx := HenryEstimate(K, x)
	return math.Abs(exact-approx) / approx
}
