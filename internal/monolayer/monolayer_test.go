package monolayer

import (
	"math"
	"testing"
)

func TestOccupiedVacant(t *testing.T) {
	if OccupiedFraction(0.3) != 0.3 {
		t.Fatal("occupied = theta")
	}
	if VacantFraction(0.3) != 0.7 {
		t.Fatal("vacant = 1 - theta")
	}
}

func TestSurfaceDensity(t *testing.T) {
	if math.Abs(SurfaceDensity(1e18, 0.5)-0.5e18) > 1e9 {
		t.Fatal("density = siteDensity * theta")
	}
}

func TestRateEquilibrium(t *testing.T) {
	// equilibrium: ka*p*(1-theta) == kd*theta for theta = K p/(1+K p)
	ka, kd, p := 2.0, 1.0, 3.0
	theta, err := EquilibriumTheta(ka, kd, p)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	r := RateIdealAdsorption(ka, kd, p, theta)
	if math.Abs(r) > 1e-9 {
		t.Fatalf("rate at equilibrium should be 0, got %v", r)
	}
	if math.Abs(theta-0.857142857) > 1e-6 {
		t.Fatalf("equilibrium theta mismatch: %v", theta)
	}
}

func TestEquilibriumThetaRejects(t *testing.T) {
	if _, err := EquilibriumTheta(0, 1, 1); err == nil {
		t.Fatal("zero ka must error")
	}
	if _, err := EquilibriumTheta(1, 1, -1); err == nil {
		t.Fatal("negative p must error")
	}
}

func TestIsostericHeatCoverageIndependent(t *testing.T) {
	if IsostericHeat(40.0, 0.1) != 40.0 {
		t.Fatal("ideal heat is coverage independent")
	}
}

func TestStepwiseLoading(t *testing.T) {
	out := StepwiseLoading(2.0, 10.0, 0.5, 3)
	if len(out) != 3 {
		t.Fatal("three steps")
	}
	// monotonic increasing toward 10
	if out[2] <= out[1] || out[1] <= out[0] {
		t.Fatal("loading must increase")
	}
	if out[2] >= 10 {
		t.Fatal("never exceeds qmax")
	}
}
