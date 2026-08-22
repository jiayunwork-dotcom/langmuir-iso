// Package comp implements the binary (two-component) Langmuir competitive
// adsorption isotherm, the natural generalisation of the single-component model
// in package iso to a surface shared by two adsorbates that compete for the
// same sites.
//
// With two adsorbates i in {1, 2}, each characterised by an equilibrium
// constant Ki and a partial pressure pi, the fractional coverage of component i
// is
//
//	theta_i = Ki pi / (1 + K1 p1 + K2 p2).
//
// The shared denominator (1 + sum over components of Ki pi) is what makes the
// model "competitive": adding more of one adsorbate displaces the other because
// the total occupancy
//
//	theta_1 + theta_2 = (K1 p1 + K2 p2) / (1 + K1 p1 + K2 p2)
//
// is always strictly below 1 (every site is either covered by one of the two
// species or empty). Setting the pressure of one component to zero recovers the
// single-component formula exactly, which is checked by the tests.
package comp
