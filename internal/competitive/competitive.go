package competitive

import (
	"errors"
)

// ErrNegativeSum is returned when the denominator of the competitive Langmuir
// model becomes non-positive, which happens only for invalid (negative) driving
// values. Under valid inputs the denominator is 1 + K1 x1 + K2 x2 and is at
// least 1, so this error marks an upstream validation failure.
var ErrNegativeSum = errors.New("competitive denominator non-positive (bad inputs)")

// Component is one adsorbate in a binary competitive Langmuir system. Both
// constants must be positive; the driving value (pressure or concentration)
// must be non-negative.
type Component struct {
	Name string
	K    float64
	Qmax float64
	X    float64
}

// Valid reports whether the component carries physically admissible constants.
func (c Component) Valid() bool {
	return c.K > 0 && c.Qmax > 0 && c.X >= 0
}

// Coverages returns the two fractional coverages for a binary mixture:
//
//	theta_i = Ki xi / (1 + K1 x1 + K2 x2).
//
// The two coverages sum to at most 1; they reach exactly 1 only when the total
// driving value diverges for both components. With one component set to zero the
// other reduces to the single-component Langmuir form.
func Coverages(a, b Component) (float64, float64, error) {
	if !a.Valid() || !b.Valid() {
		return 0, 0, ErrNegativeSum
	}
	denom := 1 + a.K*a.X + b.K*b.X
	if denom <= 0 {
		return 0, 0, ErrNegativeSum
	}
	ta := a.K * a.X / denom
	tb := b.K * b.X / denom
	return ta, tb, nil
}

// RelativeLoads returns the absolute adsorbed amounts q_i = theta_i * qmax_i.
// It is the companion of Coverages for callers that need capacities instead of
// fractions.
func RelativeLoads(a, b Component) (float64, float64, error) {
	ta, tb, err := Coverages(a, b)
	if err != nil {
		return 0, 0, err
	}
	return ta * a.Qmax, tb * b.Qmax, nil
}

// TotalCoverage is the sum of the two coverages. For a correct competitive
// model it is bounded by [0, 1]; a value above 1 signals that the denominator
// was constructed with the wrong sign somewhere upstream.
func TotalCoverage(a, b Component) (float64, error) {
	ta, tb, err := Coverages(a, b)
	if err != nil {
		return 0, err
	}
	return ta + tb, nil
}
