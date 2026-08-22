// Package inverse solves the single-component Langmuir isotherm backwards:
// given a desired fractional coverage theta it returns the driving value x that
// produces it.
//
// Inverting
//
//	theta = K x / (1 + K x)
//
// for x gives
//
//	x = theta / (K (1 - theta) ) ... no.
//
// Solving correctly:
//
//	theta (1 + K x) = K x
//	theta + theta K x = K x
//	theta = K x (1 - theta)
//	x = theta / (K (1 - theta)).
//
// The inverse is only defined for theta in [0, 1). At theta = 1 the surface is
// fully covered and no finite driving value reaches it (the isotherm only
// asymptotes to unity), so the solver returns an error rather than Infinity.
// Theta values outside [0, 1) are also rejected because coverage is a physical
// fraction of a monolayer.
package inverse
