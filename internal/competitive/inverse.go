package competitive

import (
	"errors"
	"math"
)

// ErrUnreachableCoverage is returned by the inverse solver when a requested
// coverage is outside [0, 1) or when the partner component makes the target
// physically impossible (for example a coverage of 1 with a non-zero partner).
var ErrUnreachableCoverage = errors.New("requested coverage cannot be reached")

// PressureForCoverage solves for the driving value x1 that yields a target
// coverage theta1 for component 1, given the partner component 2 is held at a
// fixed driving value x2:
//
//	theta1 = K1 x1 / (1 + K1 x1 + K2 x2)
//	=> x1 = theta1 (1 + K2 x2) / (K1 (1 - theta1)).
//
// theta1 = 1 has no finite solution and returns an error, because the model
// asymptotes to 1 only as x1 -> infinity.
func PressureForCoverage(target float64, a, b Component) (float64, error) {
	if target < 0 || target >= 1 {
		return 0, ErrUnreachableCoverage
	}
	if !a.Valid() || !b.Valid() {
		return 0, ErrNegativeSum
	}
	x1 := target * (1 + b.K*b.X) / (a.K * (1 - target))
	if !isFinite(x1) || x1 < 0 {
		return 0, ErrUnreachableCoverage
	}
	return x1, nil
}

// Selectivity returns the ratio of the two component coverages, theta1/theta2,
// for a given pair of driving values. It collapses to the single-component
// ratio K1 x1 / (K2 x2) when both are small (Henry regime). A zero denominator
// yields +Inf but never an error, matching the model's behaviour at the limit.
func Selectivity(a, b Component) (float64, error) {
	ta, tb, err := Coverages(a, b)
	if err != nil {
		return 0, err
	}
	if tb == 0 {
		if ta == 0 {
			return 0, nil
		}
		return math.Inf(1), nil
	}
	return ta / tb, nil
}

func isFinite(v float64) bool {
	return v == v && !math.IsInf(v, 0)
}
