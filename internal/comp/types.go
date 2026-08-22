package comp

// Component is the per-species description of one adsorbate in a binary system.
type Component struct {
	K    float64 // equilibrium constant (> 0)
	Qmax float64 // monolayer capacity of this species (> 0)
	P    float64 // partial pressure (>= 0)
}

// System is a binary competitive adsorption problem: two components sharing a
// homogeneous surface.
type System struct {
	A Component
	B Component
}

// CoverageSet holds the two fractional coverages and the combined total.
type CoverageSet struct {
	ThetaA float64
	ThetaB float64
	Total  float64 // ThetaA + ThetaB, always in [0, 1)
}

// AmountSet holds the adsorbed amount of each species and the total.
type AmountSet struct {
	QA float64
	QB float64
	QT float64
}

// Point is one sample of the binary isotherm at a given pressure of component A
// (component B is held fixed, which is the usual way to draw a competitive
// isotherm: sweep one adsorbate while the other sits at a constant partial
// pressure).
type Point struct {
	PA    float64
	ThetaA float64
	ThetaB float64
	Total  float64
	QA     float64
	QB     float64
}
