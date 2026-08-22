// Package iso implements the single-component Langmuir monolayer adsorption
// isotherm together with the cross-checks that any correct implementation must
// satisfy.
//
// The Langmuir model assumes a single, immobile layer of adsorbate on a
// homogeneous surface with no interactions between adsorbed molecules and no
// multilayer build-up. Under that picture the fractional surface coverage
// theta obeys
//
//	theta = K x / (1 + K x)
//
// where K is the equilibrium (adsorption) constant and x is the driving
// variable: a partial pressure p, or equivalently a bulk concentration c in
// matching units. The adsorbed amount q is the coverage multiplied by the
// monolayer capacity qmax:
//
//	q = theta * qmax.
//
// Two limiting behaviours are baked into the formula and are exercised by the
// tests in this package:
//
//   - Henry region (x -> 0): theta ~= K x, so the initial slope equals K.
//   - Saturation (x -> +inf): theta -> 1 and therefore q -> qmax.
//
// The half-coverage driving value x_{1/2} = 1/K is the point at which
// theta = 0.5. It is the single most diagnostic number of the model because it
// pins the denominator to the "1 + Kx" form: any alternative that reaches 0.5
// at a different x is a different isotherm.
//
// The package never mixes in competing models such as the Freundlich power law
// or an exponential saturating form; those are explicitly out of scope and are
// tested against in the caller's difficulty window.
package iso
