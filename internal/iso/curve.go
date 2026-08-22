package iso

import (
	"errors"
	"fmt"
	"math"
)

// ErrBadScan is returned when a Scan cannot produce any usable points (for
// example an empty explicit list, a non-positive count, or a zero/negative
// span for a linear scan).
var ErrBadScan = errors.New("invalid scan specification")

// Points generates the isotherm samples described by s for the given constants.
// If s.Press... no, if s.Pressures is non-empty those exact driving values are
// used. Otherwise N points are laid out between XMin and XMax.
//
// For a Linear scan the points are x_i = XMin + i*(XMax-XMin)/(N-1). For a Log
// scan the points are 10^(log10(XMin) + i*(log10(XMax)-log10(XMin))/(N-1)),
// which requires XMin > 0. Each point records its coverage and adsorbed amount.
//
// The coverage is clamped to [0, 1) only by virtue of the model itself; the
// function does not post-process results, so a caller that supplies K <= 0 or
// negative pressures will see the (non-physical) raw formula output and should
// have validated first.
func Points(K, qmax float64, s Scan) ([]Point, error) {
	if err := Validate(K, qmax, 0); err != nil {
		return nil, err
	}
	xs, err := sampleDriving(s)
	if err != nil {
		return nil, err
	}
	out := make([]Point, 0, len(xs))
	for _, x := range xs {
		if x < 0 {
			return nil, fmt.Errorf("%w: negative driving value %v produced by scan", ErrBadScan, x)
		}
		theta := Coverage(K, x)
		out = append(out, Point{X: x, Theta: theta, Q: Adsorption(theta, qmax)})
	}
	return out, nil
}

// sampleDriving turns a Scan into the concrete list of driving values.
func sampleDriving(s Scan) ([]float64, error) {
	if len(s.Pressures) > 0 {
		xs := make([]float64, len(s.Pressures))
		copy(xs, s.Pressures)
		return xs, nil
	}
	if s.N <= 0 {
		return nil, fmt.Errorf("%w: n must be > 0 (got %d)", ErrBadScan, s.N)
	}
	if s.XMax <= s.XMin {
		return nil, fmt.Errorf("%w: xmax (%v) must exceed xmin (%v)", ErrBadScan, s.XMax, s.XMin)
	}
	switch s.Scale {
	case Linear:
		return linearSpace(s.XMin, s.XMax, s.N), nil
	case Log:
		if s.XMin <= 0 {
			return nil, fmt.Errorf("%w: log scale requires xmin > 0 (got %v)", ErrBadScan, s.XMin)
		}
		return logSpace(s.XMin, s.XMax, s.N), nil
	default:
		return nil, fmt.Errorf("%w: unknown scale", ErrBadScan)
	}
}

func linearSpace(lo, hi float64, n int) []float64 {
	xs := make([]float64, n)
	if n == 1 {
		xs[0] = lo
		return xs
	}
	step := (hi - lo) / float64(n-1)
	for i := 0; i < n; i++ {
		xs[i] = lo + float64(i)*step
	}
	return xs
}

func logSpace(lo, hi float64, n int) []float64 {
	llo := math.Log10(lo)
	lhi := math.Log10(hi)
	xs := make([]float64, n)
	if n == 1 {
		xs[0] = lo
		return xs
	}
	step := (lhi - llo) / float64(n-1)
	for i := 0; i < n; i++ {
		xs[i] = math.Pow(10, llo+float64(i)*step)
	}
	return xs
}

// MaxTheta returns the largest coverage in a slice of points. Because the model
// is monotonic in x, this is the value at the largest driving value and is the
// closest approach to saturation present in the sample.
func MaxTheta(pts []Point) float64 {
	max := 0.0
	for _, p := range pts {
		if p.Theta > max {
			max = p.Theta
		}
	}
	return max
}
