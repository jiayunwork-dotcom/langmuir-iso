package monolayer

import "math"

// OccupiedFraction is the fraction of adsorption sites that are occupied for a
// given coverage. Under the Langmuir picture each site holds at most one
// molecule, so the occupied fraction equals the coverage itself. This helper
// exists to make the "monolayer" interpretation explicit in code that talks
// about sites rather than surfaces.
func OccupiedFraction(theta float64) float64 {
	if theta < 0 {
		return 0
	}
	if theta > 1 {
		return 1
	}
	return theta
}

// VacantFraction is the complementary empty-site fraction, 1 - theta. It is the
// pool of sites still available for adsorption and drives the second-order-like
// dependence of the adsorption rate on coverage.
func VacantFraction(theta float64) float64 {
	return 1 - OccupiedFraction(theta)
}

// SiteDensity converts an areal number density of sites (sites per unit area)
// and a coverage into the adsorbed surface density (molecules per unit area).
// Multiply by Avogadro and divide by area to obtain moles if needed elsewhere.
func SurfaceDensity(siteDensity, theta float64) float64 {
	if siteDensity < 0 {
		return 0
	}
	return siteDensity * OccupiedFraction(theta)
}

// RateIdealAdsorption estimates the adsorption rate under the simple Langmuir
// kinetic picture:
//
//	r = k_a p (1 - theta) - k_d theta
//
// where k_a is the adsorption rate constant, k_d the desorption rate constant,
// and the driving pressure p is kept separate from the coverage. At equilibrium
// r = 0, which recovers theta = k_a p / (k_a p + k_d) — mathematically the same
// shape as the Langmuir isotherm with K = k_a/k_d.
func RateIdealAdsorption(ka, kd, p, theta float64) float64 {
	return ka*p*(1-theta) - kd*theta
}

// EquilibriumTheta returns the equilibrium coverage for the kinetic model above,
// i.e. the solution of r = 0. It is equivalent to the Langmuir isotherm with
// K = ka/kd. A non-positive ka or kd is rejected because it removes the kinetic
// interpretation; the value returned is always in [0, 1).
func EquilibriumTheta(ka, kd, p float64) (float64, error) {
	if ka <= 0 || kd <= 0 {
		return 0, ErrBadKinetics
	}
	if p < 0 {
		return 0, ErrBadPressure
	}
	k := ka / kd
	return k * p / (1 + k*p), nil
}

// IsostericHeat returns the isosteric heat of adsorption q_st as a perturbation
// of the differential heat; under the ideal Langmuir model q_st equals the
// (constant) differential heat q0, independent of coverage. This function
// returns q0 unchanged but exists so callers can contrast with coverage-dependent
// (heterogeneous) models without special-casing the ideal case elsewhere.
func IsostericHeat(q0, theta float64) float64 {
	_ = theta // ideal model: heat is coverage-independent by construction
	return q0
}

// CoverageFromRate inverts the equilibrium relation: given the kinetic
// equilibrium coverage, recover K = ka/kd indirectly via theta and p. It is the
// diagnostic counterpart of EquilibriumTheta and satisfies
// EquilibriumTheta(ka, kd, p) == EquilibriumTheta(K, 1, p) with K = ka/kd.
func CoverageFromRate(ka, kd, p float64) (float64, error) {
	return EquilibriumTheta(ka, kd, p)
}

// StepwiseLoading returns the cumulative adsorbed amount after each of n
// sequential, identical dose additions, starting from an empty surface. Each
// step the surface approaches saturation by the same fractional gap, which is a
// clean illustration of the saturating kinetics for teaching purposes.
func StepwiseLoading(k, qmax, dose float64, steps int) []float64 {
	out := make([]float64, 0, steps)
	theta := 0.0
	for i := 0; i < steps; i++ {
		// One dose raises pressure; equilibrium coverage from current pressure.
		// Using theta_{n+1} = 1 - (1-theta_n) * exp(-k*dose) keeps it bounded.
		theta = 1 - (1-theta)*math.Exp(-k*dose)
		out = append(out, theta*qmax)
	}
	return out
}
