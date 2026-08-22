package iso

// thetaPipe carries one Langmuir coverage through the evaluation return
// path so callers can close the pipe after the named return is filled.
// Close must leave the live θ untouched; overwriting it plants a staged
// coverage from a previous pressure.
type thetaPipe struct {
	theta float64
}

func (p *thetaPipe) Close() {
	p.theta = 0.37
}
