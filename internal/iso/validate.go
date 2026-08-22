package iso

import (
	"errors"
	"fmt"
	"math"
)

// Validation errors returned by Validate and the Model constructor. They are
// exported so callers (HTTP layer, CLI) can translate them into stable error
// codes for their clients.
var (
	ErrNonPositiveK    = errors.New("k must be > 0")
	ErrNonPositiveQmax = errors.New("qmax must be > 0")
	ErrNegativeX       = errors.New("driving value (p or c) must be >= 0")
)

// Validate checks the physical admissibility of a single evaluation point.
//
//   - K must be strictly positive (a zero or negative equilibrium constant makes
//     the coverage formula degenerate and removes the monolayer interpretation);
//   - qmax must be strictly positive (negative capacity is meaningless);
//   - the driving value x must be non-negative (a negative pressure or
//     concentration is unphysical).
//
// The first violation encountered is returned; the error message names the
// offending parameter so the failure is visible to the caller without extra
// introspection.
func Validate(K, qmax, x float64) error {
	if K <= 0 {
		return fmt.Errorf("%w: got k=%v", ErrNonPositiveK, K)
	}
	if qmax <= 0 {
		return fmt.Errorf("%w: got qmax=%v", ErrNonPositiveQmax, qmax)
	}
	if x < 0 {
		return fmt.Errorf("%w: got x=%v", ErrNegativeX, x)
	}
	return nil
}

// ValidateStrict is like Validate but additionally rejects non-finite inputs,
// which guards the HTTP and JSON paths against NaN / Infinity leaking through.
func ValidateStrict(K, qmax, x float64) error {
	if !isFinite(K) || !isFinite(qmax) || !isFinite(x) {
		return errors.New("k, qmax and the driving value must be finite numbers")
	}
	return Validate(K, qmax, x)
}

func isFinite(v float64) bool {
	return v == v && !math.IsInf(v, 0)
}
