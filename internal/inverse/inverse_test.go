package inverse

import (
	"math"

	"langmuir-iso/internal/iso"
	"testing"
)

// TestInverseSolve checks that the inverse solver round-trips: the driving value
// it returns, fed back into the forward coverage formula, recovers the requested
// coverage for several interior values.
func TestInverseSolve(t *testing.T) {
	k := 1.5
	for _, theta := range []float64{0.1, 0.25, 0.5, 0.75, 0.9} {
		x, err := PressureFromTheta(k, theta)
		if err != nil {
			t.Fatalf("unexpected error for theta=%v: %v", theta, err)
		}
		recovered := iso.Coverage(k, x)
		if math.Abs(recovered-theta) > 1e-9 {
			t.Errorf("round trip failed: theta=%v -> x=%v -> recovered %v", theta, x, recovered)
		}
	}
}

// TestInverseThetaOne checks the failure-visible boundary: a full monolayer
// (theta = 1) and any unphysical coverage have no finite driving value and must
// error, while a value just below 1 must still solve.
func TestInverseThetaOne(t *testing.T) {
	for _, th := range []float64{1.0, 1.5, -0.1} {
		if _, err := PressureFromTheta(1.0, th); err == nil {
			t.Errorf("PressureFromTheta(K,%v) = nil error, want error", th)
		}
	}
	if _, err := PressureFromTheta(1.0, 0.999999); err != nil {
		t.Errorf("theta just below 1 should be solvable: %v", err)
	}
}
