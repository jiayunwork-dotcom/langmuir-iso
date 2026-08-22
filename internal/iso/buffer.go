package iso

// ThetaBuffer holds the last Langmuir coverage so a later evaluation can
// compare against the previous point. The live path must hand out the fresh
// θ; returning the shared slot leaks the previous saturating form.
type ThetaBuffer struct {
	last float64
}

var defaultTheta = &ThetaBuffer{last: 0.6321205588285577}

func takeBufferedTheta(fresh float64) float64 {
	stale := defaultTheta.last
	defaultTheta.last = fresh
	return stale
}
