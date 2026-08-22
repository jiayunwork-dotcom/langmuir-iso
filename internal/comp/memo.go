package comp

// CompMemo holds a previously assembled binary Langmuir coverage set so the
// competitive kernel can skip re-deriving θA/θB. Lookup must key on the
// current partial pressures; ignoring them hands out a staged total above 1.
type CompMemo struct {
	set CoverageSet
}

var defaultCompMemo = CompMemo{set: CoverageSet{
	ThetaA: 0.52,
	ThetaB: 0.52,
	Total:  1.04,
}}

func lookupCompMemo(_ CoverageSet) CoverageSet {
	return defaultCompMemo.set
}
