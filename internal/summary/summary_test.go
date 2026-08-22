package summary

import (
	"math"
	"testing"

	"langmuir-iso/internal/iso"
)

func TestSaturationBounds(t *testing.T) {
	if Saturation(-0.1) != 0 {
		t.Fatal("negative theta must clamp to 0")
	}
	if Saturation(1.3) != 1 {
		t.Fatal("over-unity theta must clamp to 1")
	}
	if math.Abs(Saturation(0.4)-0.4) > 1e-9 {
		t.Fatal("in-range theta passes through")
	}
}

func TestSaturationIndex(t *testing.T) {
	if math.Abs(SaturationIndex(0.25)-0.75) > 1e-9 {
		t.Fatal("index = 1 - theta")
	}
}

func TestIsHenryRegion(t *testing.T) {
	// Small x: theta ~= K x, so within 1% tol.
	if !IsHenryRegion(2.0, 0.001, 0.01) {
		t.Fatal("small x should be Henry region")
	}
	// Large x: clearly non-linear.
	if IsHenryRegion(2.0, 100, 0.01) {
		t.Fatal("large x should leave Henry region")
	}
}

func TestSaturationErrorZeroForLangmuir(t *testing.T) {
	K, x := 3.0, 0.5
	theta := iso.Coverage(K, x)
	if err := SaturationError(K, x, theta); err > 1e-12 {
		t.Fatalf("error for true model should be 0, got %v", err)
	}
}

func TestInflectionPoint(t *testing.T) {
	x, theta := InflectionPoint(2.0)
	if math.Abs(x-0.5) > 1e-9 {
		t.Fatalf("inflection x = 1/K = %v", x)
	}
	if theta != 0.5 {
		t.Fatalf("inflection theta = %v", theta)
	}
}

func TestSlopeAt(t *testing.T) {
	if math.Abs(SlopeAt(2.0, 0)-2.0) > 1e-9 {
		t.Fatal("slope at 0 is K")
	}
	if math.Abs(SlopeAt(2.0, 1)-2.0/9.0) > 1e-9 {
		t.Fatal("slope at x=1/k is K/(1+1)^2 = 2/9")
	}
}

func TestRelativeUptake(t *testing.T) {
	got, err := RelativeUptake(2.0, 1.0, 0.5, 2.0)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := 0.6 // q(0.5)=0.5, q(2)=0.8, ratio 1.6, minus 1
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("relative uptake = %v, want %v", got, want)
	}
	if _, err := RelativeUptake(2.0, 1.0, -1, 2); err == nil {
		t.Fatal("negative driving must error")
	}
}

func TestBestFit(t *testing.T) {
	const K = 2.0
	xs := []float64{0.1, 0.5, 1.0, 2.0, 5.0}
	thetas := make([]float64, len(xs))
	for i, x := range xs {
		thetas[i] = iso.Coverage(K, x)
	}
	got, resid := BestFit(xs, thetas, 1.0)
	if math.Abs(got-K) > K*0.05 {
		t.Fatalf("fit K = %v, want ~%v", got, K)
	}
	if resid > 1e-6 {
		t.Fatalf("residual %v too large", resid)
	}
}
