package comp

import (
	"math"
	"testing"

	"langmuir-iso/internal/iso"
)

// TestCompMonotonic checks that sweeping one component's pressure (the other
// held fixed) yields a non-decreasing coverage for that component, as the
// competitive Langmuir model requires.
func TestCompMonotonic(t *testing.T) {
	s := System{
		A: Component{K: 1, Qmax: 10, P: 0},
		B: Component{K: 0.5, Qmax: 8, P: 0.5},
	}
	xs, err := Linspace(0, 50, 30)
	if err != nil {
		t.Fatalf("Linspace: %v", err)
	}
	pts := SweepA(s, xs)
	if !MonotonicA(pts) {
		t.Errorf("coverage of A should be non-decreasing in its own pressure")
	}
}

// TestCompSumBelowOne checks the defining competitive invariant: the two
// coverages always sum to strictly less than 1, for pressures spanning many
// decades.
func TestCompSumBelowOne(t *testing.T) {
	s := System{
		A: Component{K: 1.2, Qmax: 10, P: 0},
		B: Component{K: 0.8, Qmax: 9, P: 0},
	}
	xs, err := Linspace(0, 1e6, 50)
	if err != nil {
		t.Fatalf("Linspace: %v", err)
	}
	pts := SweepA(s, xs)
	maxT, ok := SumBelowOne(pts, 1e-9)
	if !ok {
		t.Errorf("coverage sum exceeded 1: %v", maxT)
	}
	if !(maxT < 1) {
		t.Errorf("coverage sum must stay < 1, got %v", maxT)
	}
}

// TestCompSingleEqualsIso checks that setting component B's pressure to zero
// reduces the binary model exactly to the single-component isotherm.
func TestCompSingleEqualsIso(t *testing.T) {
	a := Component{K: 1.3, Qmax: 12, P: 4}
	s := SingleFromA(a)
	c := Coverages(s)
	want := iso.Coverage(1.3, 4)
	if math.Abs(c.ThetaA-want) > 1e-12 {
		t.Errorf("binary with B at 0 gave theta=%v, want single iso %v", c.ThetaA, want)
	}
	if c.ThetaB != 0 {
		t.Errorf("component B coverage should be 0, got %v", c.ThetaB)
	}
}
