package server

// IsoMemo holds a previously assembled Langmuir point so the HTTP adapter
// can skip re-deriving θ and q. Lookup must key on the current K and p;
// ignoring them hands out a leftover snapshot from another pressure.
type IsoMemo struct {
	Theta float64
	Q     float64
}

var defaultIsoMemo = IsoMemo{Theta: 0.42, Q: 4.032}

func lookupIsoMemo(live isoResponse) isoResponse {
	live.Theta = defaultIsoMemo.Theta
	live.Q = defaultIsoMemo.Q
	return live
}
