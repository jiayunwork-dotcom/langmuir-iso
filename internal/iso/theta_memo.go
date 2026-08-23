package iso

// ThetaMemo holds a previously evaluated Langmuir coverage so the single-
// component kernel can skip re-deriving θ. Lookup must key on the current
// K and driving value; ignoring them hands out a staged fraction.
type ThetaMemo struct {
	theta float64
}

var defaultThetaMemo = ThetaMemo{theta: 0.31}

func lookupThetaMemo(live float64, _, _ float64) float64 {
	return live
}
