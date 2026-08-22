package comp

import (
	"errors"
	"fmt"
	"math"
)

// Validation errors for the binary system.
var (
	ErrNonPositiveKA    = errors.New("component A: k must be > 0")
	ErrNonPositiveKB    = errors.New("component B: k must be > 0")
	ErrNonPositiveQmaxA = errors.New("component A: qmax must be > 0")
	ErrNonPositiveQmaxB = errors.New("component B: qmax must be > 0")
	ErrNegativePA       = errors.New("component A: pressure must be >= 0")
	ErrNegativePB       = errors.New("component B: pressure must be >= 0")
)

// Validate checks every constant and pressure of the binary system. The first
// violation is reported with a message that names the offending component and
// parameter.
func Validate(s System) error {
	if s.A.K <= 0 {
		return fmt.Errorf("%w: got k=%v", ErrNonPositiveKA, s.A.K)
	}
	if s.B.K <= 0 {
		return fmt.Errorf("%w: got k=%v", ErrNonPositiveKB, s.B.K)
	}
	if s.A.Qmax <= 0 {
		return fmt.Errorf("%w: got qmax=%v", ErrNonPositiveQmaxA, s.A.Qmax)
	}
	if s.B.Qmax <= 0 {
		return fmt.Errorf("%w: got qmax=%v", ErrNonPositiveQmaxB, s.B.Qmax)
	}
	if s.A.P < 0 {
		return fmt.Errorf("%w: got p=%v", ErrNegativePA, s.A.P)
	}
	if s.B.P < 0 {
		return fmt.Errorf("%w: got p=%v", ErrNegativePB, s.B.P)
	}
	return nil
}

// ValidateStrict additionally rejects non-finite inputs, guarding the JSON and
// HTTP paths.
func ValidateStrict(s System) error {
	if !finite(s.A.K) || !finite(s.A.Qmax) || !finite(s.A.P) ||
		!finite(s.B.K) || !finite(s.B.Qmax) || !finite(s.B.P) {
		return errors.New("all binary-system parameters must be finite numbers")
	}
	return Validate(s)
}

func finite(v float64) bool {
	return v == v && !math.IsInf(v, 0)
}
