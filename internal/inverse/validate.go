package inverse

import (
	"errors"
	"fmt"
)

// Validation errors for the inverse problem.
var (
	ErrNonPositiveKInv = errors.New("inverse: k must be > 0")
	ErrThetaRange      = errors.New("inverse: theta must be in [0, 1)")
)

// Validate checks the structured inverse problem. It is a thin, named wrapper
// around PressureFromTheta's preconditions so the HTTP layer can produce a
// stable error code without recomputing the message.
func Validate(p Problem) error {
	if !isFinite(p.K) {
		return fmt.Errorf("%w: got k=%v", ErrNonPositiveKInv, p.K)
	}
	if p.K <= 0 {
		return fmt.Errorf("%w: got k=%v", ErrNonPositiveKInv, p.K)
	}
	if !isFinite(p.Theta) {
		return fmt.Errorf("%w: got theta=%v", ErrThetaRange, p.Theta)
	}
	if p.Theta < 0 || p.Theta >= 1 {
		return fmt.Errorf("%w: got theta=%v", ErrThetaRange, p.Theta)
	}
	return nil
}
