package competitive

import (
	"math"
	"testing"
)

func TestCoveragesSingleComponent(t *testing.T) {
	a := Component{Name: "A", K: 1.0, Qmax: 10, X: 1.0}
	b := Component{Name: "B", K: 1.0, Qmax: 10, X: 0}
	ta, tb, err := Coverages(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// B is absent, so A must reduce to single-component Langmuir with K*x/(1+Kx).
	want := 1.0 * 1.0 / (1 + 1.0*1.0)
	if math.Abs(ta-want) > 1e-9 {
		t.Fatalf("ta = %v, want %v", ta, want)
	}
	if tb != 0 {
		t.Fatalf("tb = %v, want 0", tb)
	}
}

func TestCoveragesSumAtMostOne(t *testing.T) {
	a := Component{Name: "A", K: 2.0, Qmax: 5, X: 3.0}
	b := Component{Name: "B", K: 0.5, Qmax: 5, X: 4.0}
	sum, err := TotalCoverage(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum > 1 {
		t.Fatalf("total coverage %v exceeds 1", sum)
	}
}

func TestCoveragesInvalid(t *testing.T) {
	a := Component{Name: "A", K: -1, Qmax: 5, X: 1.0}
	b := Component{Name: "B", K: 1.0, Qmax: 5, X: 1.0}
	if _, _, err := Coverages(a, b); err == nil {
		t.Fatal("expected error for negative K")
	}
}

func TestPressureForCoverage(t *testing.T) {
	a := Component{Name: "A", K: 1.0, Qmax: 10, X: 0}
	b := Component{Name: "B", K: 1.0, Qmax: 10, X: 1.0}
	x1, err := PressureForCoverage(0.5, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// x1 = 0.5*(1 + 1*1)/(1*(1-0.5)) = 0.5*2/0.5 = 2
	if math.Abs(x1-2) > 1e-9 {
		t.Fatalf("x1 = %v, want 2", x1)
	}
}

func TestPressureForCoverageUnityRejected(t *testing.T) {
	a := Component{Name: "A", K: 1.0, Qmax: 10, X: 0}
	b := Component{Name: "B", K: 1.0, Qmax: 10, X: 1.0}
	if _, err := PressureForCoverage(1.0, a, b); err == nil {
		t.Fatal("theta=1 must be unreachable")
	}
}

func TestSelectivityHenry(t *testing.T) {
	a := Component{Name: "A", K: 2.0, Qmax: 5, X: 0.01}
	b := Component{Name: "B", K: 1.0, Qmax: 5, X: 0.02}
	s, err := Selectivity(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Henry limit: K1 x1 / (K2 x2) = 2*0.01 / (1*0.02) = 1
	if math.Abs(s-1) > 1e-6 {
		t.Fatalf("selectivity = %v, want ~1", s)
	}
}
