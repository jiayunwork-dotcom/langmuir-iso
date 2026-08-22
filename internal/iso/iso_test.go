package iso

import (
	"math"
	"testing"
)

// TestValidateNegative checks the three hard validation rules: K must be
// positive, qmax must be positive and the driving value must be non-negative.
// These are the simple, single-rule runtime guards whose absence would let a
// degenerate input slip through silently.
func TestValidateNegative(t *testing.T) {
	bad := []struct {
		k, qmax, x float64
	}{
		{0, 10, 1},   // K = 0
		{-1, 10, 1},  // K < 0
		{1, 0, 1},    // qmax = 0
		{1, -5, 1},   // qmax < 0
		{1, 10, -1},  // p < 0
	}
	for _, c := range bad {
		if err := Validate(c.k, c.qmax, c.x); err == nil {
			t.Errorf("Validate(%v,%v,%v) = nil, want error", c.k, c.qmax, c.x)
		}
	}
	if err := Validate(1, 10, 0); err != nil {
		t.Errorf("Validate(1,10,0) = %v, want nil", err)
	}
}

// TestHalfPressure checks the analytic half-coverage pressure x_{1/2} = 1/K.
func TestHalfPressure(t *testing.T) {
	for _, k := range []float64{0.5, 1, 2, 3.7} {
		got := HalfPressure(k)
		want := 1 / k
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("HalfPressure(%v) = %v, want %v", k, got, want)
		}
	}
}

// TestHenryLimit checks that at small driving values theta/x approaches the
// Henry constant K (the initial slope of the isotherm).
func TestHenryLimit(t *testing.T) {
	k := 2.0
	for _, x := range []float64{1e-4, 1e-5, 1e-6} {
		theta := Coverage(k, x)
		ratio := theta / x
		if math.Abs(ratio-k) > 1e-2 {
			t.Errorf("theta/x at x=%v = %v, want ~K=%v", x, ratio, k)
		}
	}
	if got, want := HenrySlope(k), k; got != want {
		t.Errorf("HenrySlope(%v) = %v, want %v", k, got, want)
	}
}

// TestCoverageAtInfinity checks the two saturation limits: at p=0 the coverage is
// exactly 0, and as p grows without bound the coverage approaches 1 from below
// and never exceeds it.
func TestCoverageAtInfinity(t *testing.T) {
	k := 1.0
	if Coverage(k, 0) != 0 {
		t.Errorf("Coverage(K,0) = %v, want 0", Coverage(k, 0))
	}
	theta := Coverage(k, 1e12)
	if theta >= 1 {
		t.Errorf("Coverage at large x = %v, must stay < 1", theta)
	}
	if math.Abs(theta-1) > 1e-9 {
		t.Errorf("Coverage at large x = %v, want ~1", theta)
	}
}

// TestCrossDoubleK checks the cross-rule that doubling K at a fixed driving
// value raises the coverage, and halves the half-coverage pressure.
func TestCrossDoubleK(t *testing.T) {
	p := 1.0
	t1 := Coverage(1, p)
	t2 := Coverage(2, p)
	if !(t2 > t1) {
		t.Errorf("doubling K at fixed p: theta %v not > %v", t2, t1)
	}
	xh1 := HalfPressure(1)
	xh2 := HalfPressure(2)
	if math.Abs(xh2-xh1/2) > 1e-12 {
		t.Errorf("doubling K should halve x_{1/2}: got %v want %v", xh2, xh1/2)
	}
}

// TestCrossDoubleQmax checks the cross-rule that changing qmax leaves the
// coverage unchanged while doubling the adsorbed amount.
func TestCrossDoubleQmax(t *testing.T) {
	k, p := 1.0, 1.0
	theta := Coverage(k, p)
	q1 := Adsorption(theta, 10)
	q2 := Adsorption(theta, 20)
	if math.Abs(q1-10*theta) > 1e-12 || math.Abs(q2-20*theta) > 1e-12 {
		t.Errorf("q must equal qmax*theta: q1=%v q2=%v", q1, q2)
	}
	if math.Abs(q2-2*q1) > 1e-12 {
		t.Errorf("doubling qmax should double q: got %v want %v", q2, 2*q1)
	}
	if math.Abs(Coverage(k, p)-theta) > 0 {
		t.Errorf("theta must be independent of qmax")
	}
}

// TestCoverageAtHalfPressure is the defining Langmuir invariant and the hardest
// single check: evaluated at the half-coverage driving value x_{1/2} = 1/K the
// coverage must be exactly 0.5 for every K. This pins the denominator to the
// "1 + Kx" form; an alternative saturating model (for example 1 - exp(-Kx),
// which reaches 0.5 at x = ln(2)/K instead) fails this test deterministically.
func TestCoverageAtHalfPressure(t *testing.T) {
	for _, k := range []float64{0.3, 1, 2.5, 7.0} {
		xh := HalfPressure(k)
		theta := Coverage(k, xh)
		if math.Abs(theta-0.5) > 1e-9 {
			t.Errorf("theta at x_{1/2}=1/K = %v, want exactly 0.5 (K=%v)", theta, k)
		}
	}
}
