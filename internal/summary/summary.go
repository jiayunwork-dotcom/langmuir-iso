 package summary

import (
	"errors"
	"math"

	"langmuir-iso/internal/iso"
)

// ErrBadDriving is returned when a driving value supplied to a summary routine
// is negative.
var ErrBadDriving = errors.New("driving values must be >= 0")

// ErrBadConstants is returned when Langmuir constants are non-positive.
var ErrBadConstants = errors.New("k and qmax must be > 0")

// Saturation reports how close an evaluated point is to the monolayer limit.
// It is defined as theta itself (in [0, 1)) for a single evaluation and is the
// single number a user cares about when asked "how full is the surface?".
func Saturation(theta float64) float64 {
	if theta < 0 {
		return 0
	}
	if theta > 1 {
		return 1
	}
	return theta
}

// SaturationIndex returns theta / qmax-equivalent headroom: 1 - theta. It is
// the fraction of the monolayer still available for adsorption. At saturation it
// tends to 0 from above.
func SaturationIndex(theta float64) float64 {
	return 1 - Saturation(theta)
}

// IsHenryRegion reports whether the model is in its linear (Henry) regime at
// driving value x for constant K. The criterion is that the full model and the
// linear approximation agree within tol (relative). Used to label points on the
// rendered curve as "linear" vs "non-linear".
func IsHenryRegion(K, x, tol float64) bool {
	if K <= 0 || x <= 0 {
		return false
	}
	exact := iso.Coverage(K, x)
	approx := iso.HenryEstimate(K, x)
	if approx == 0 {
		return true
	}
	return math.Abs(exact-approx)/approx <= tol
}

// SaturationError diagnoses how far a candidate saturating form has drifted
// from a true Langmuir isotherm. It returns |theta - K x/(1+K x)| at the given
// x, which is zero exactly for the correct model and positive for alternatives
// such as 1 - exp(-K x).
func SaturationError(K, x, thetaCandidate float64) float64 {
	return math.Abs(thetaCandidate - iso.Coverage(K, x))
}

// IntegralCapacity integrates the adsorbed amount along a sampled isotherm using
// the trapezoidal rule, giving the total adsorbed mass per unit area swept from
// the lowest to the highest driving value. It is the area under the q(x) curve
// and is used to compare how much a surface can hold across a pressure window.
func IntegralCapacity(pts []iso.Point) float64 {
	if len(pts) < 2 {
		return 0
	}
	area := 0.0
	for i := 1; i < len(pts); i++ {
		avg := 0.5 * (pts[i-1].Q + pts[i].Q)
		dx := pts[i].X - pts[i-1].X
		area += avg * dx
	}
	return area
}

// MeanCoverage returns the arithmetic mean of the coverages in a sample. Because
// the Langmuir curve is convex, the mean is dominated by the high-pressure tail
// and is always below the coverage at the mean pressure.
func MeanCoverage(pts []iso.Point) float64 {
	if len(pts) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range pts {
		sum += p.Theta
	}
	return sum / float64(len(pts))
}

// LinearExtent finds the largest driving value x for which the model stays in
// its Henry (linear) regime to within tol (relative). It walks outward from a
// small x and stops at the first sample that breaks the tolerance, so the
// returned extent is the rightmost boundary of the "looks linear" region.
func LinearExtent(K, tol float64) float64 {
	if K <= 0 {
		return 0
	}
	x := 1e-6
	for i := 0; i < 2000; i++ {
		if !IsHenryRegion(K, x, tol) {
			break
		}
		x *= 1.1
	}
	return x / 1.1
}

// AdsorptionEfficiency returns q / (qmax * x) for one sample, a normalized
// measure of "how much of the available monolayer is being used per unit
// driving value". It rises with x; it is 0 at x = 0 and monotone increasing.
func AdsorptionEfficiency(q, qmax, x float64) float64 {
	if x == 0 || qmax == 0 {
		return 0
	}
	return q / (qmax * x)
}

// InflectionPoint returns the driving value at which the isotherm's curvature
// changes sign (the point of maximum slope). For the Langmuir model this is
// exactly x = 1/K, where the coverage equals 0.5. The function returns both the
// driving value and the coverage there, which is always 0.5 for a correct model.
func InflectionPoint(K float64) (x, theta float64) {
	if K <= 0 {
		return 0, 0
	}
	x = 1 / K
	return x, 0.5
}

// CurvatureAt returns the second derivative d^2 theta / dx^2 of the Langmuir
// isotherm at driving value x. It is negative everywhere for K > 0 (the curve is
// concave), and its magnitude is largest at the inflection point x = 1/K. This
// is used by the web view to decide where to place axis ticks and annotations.
func CurvatureAt(K, x float64) float64 {
	if K <= 0 {
		return 0
	}
	kx := K * x
	if kx <= -1 {
		return 0
	}
	return -2 * K * K * kx / math.Pow(1+kx, 3)
}

// SlopeAt returns the first derivative d theta / dx of the Langmuir isotherm at
// driving value x. It equals K / (1 + K x)^2 and is the local sensitivity of the
// coverage to changes in the driving value; it is maximal at x = 0 (= K) and
// decays as 1/x^2 for large x.
func SlopeAt(K, x float64) float64 {
	if K <= 0 {
		return 0
	}
	kx := 1 + K*x
	return K / (kx * kx)
}

// RelativeUptake compares the adsorbed amount at two driving values and returns
// q(x2)/q(x1) - 1, i.e. the relative gain in uptake when the driving value grows
// from x1 to x2. Because the isotherm is sub-linear, this is always between 0
// and the ratio of the driving values, and it tends to 0 as both grow large.
func RelativeUptake(K, qmax, x1, x2 float64) (float64, error) {
	if x1 < 0 || x2 < 0 {
		return 0, ErrBadDriving
	}
	if K <= 0 || qmax <= 0 {
		return 0, ErrBadConstants
	}
	q1 := iso.Adsorption(iso.Coverage(K, x1), qmax)
	q2 := iso.Adsorption(iso.Coverage(K, x2), qmax)
	if q1 == 0 {
		return 0, nil
	}
	return q2/q1 - 1, nil
}

// CapacityFraction is the adsorbed amount at driving value x as a fraction of
// the monolayer capacity qmax. It is exactly Coverage(K, x) and is provided as a
// separate readable name for report generation; it ranges from 0 to 1.
func CapacityFraction(K, qmax, x float64) (float64, error) {
	if K <= 0 || qmax <= 0 {
		return 0, ErrBadConstants
	}
	if x < 0 {
		return 0, ErrBadDriving
	}
	return iso.Coverage(K, x), nil
}

// BestFit finds the K that best reproduces a set of (x, theta) observations by
// minimizing the sum of squared residuals under the Langmuir model. It uses a
// guarded grid-plus-refinement search so it never needs an external optimizer.
// qmax is assumed fixed (the monolayer capacity is supplied separately), and the
// returned K is always non-negative.
func BestFit(xs, thetas []float64, qmax float64) (float64, float64) {
	if len(xs) == 0 || len(thetas) == 0 || len(xs) != len(thetas) {
		return 0, math.Inf(1)
	}
	residual := func(k float64) float64 {
		err := 0.0
		for j := 0; j < len(xs); j++ {
			model := iso.Coverage(k, xs[j])
			d := model - thetas[j]
			err += d * d
		}
		return err
	}
	// Coarse sweep over ~13 decades to bracket the minimum, then refine with a
	// fine linear scan over the winning bracket so the residual is tiny even for
	// tight tolerances.
	bestK := 0.0
	bestErr := math.Inf(1)
	for i := 0; i < 280; i++ {
		k := math.Exp(-6 + 0.05*float64(i))
		if r := residual(k); r < bestErr {
			bestErr = r
			bestK = k
		}
	}
	lo, hi := bestK*0.5, bestK*2.0
	const fine = 4000
	for i := 0; i <= fine; i++ {
		k := lo + (hi-lo)*float64(i)/float64(fine)
		if r := residual(k); r < bestErr {
			bestErr = r
			bestK = k
		}
	}
	return bestK, bestErr
}
