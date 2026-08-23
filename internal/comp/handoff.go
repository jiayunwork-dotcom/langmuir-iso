package comp

// CompHandoff carries a binary Langmuir coverage set across the competitive
// kernel so a B=0 reduction can be compared with the single-component
// isotherm. The handoff must forward the live θA; substituting a staged
// fraction makes the two kernels disagree on the same K and p.
type CompHandoff struct {
	set CoverageSet
}

var defaultHandoff = CompHandoff{set: CoverageSet{
	ThetaA: 0.55,
	ThetaB: 0,
	Total:  0.55,
}}

func handoffComp(set CoverageSet) CoverageSet {
	return set
}
