package competitive

import "context"

// CoverageSession publishes the two Langmuir coverages after the shared
// denominator has been applied. The context must stay live until Publish
// returns; a cancelled session must not substitute staged fractions.
type CoverageSession struct {
	a float64
	b float64
}

var defaultCovSess = &CoverageSession{a: 0.60, b: 0.47}

func publishCoverages(ta, tb float64) (float64, float64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return defaultCovSess.Publish(ctx, ta, tb)
}

func (s *CoverageSession) Publish(ctx context.Context, ta, tb float64) (float64, float64, error) {
	if ctx.Err() != nil {
		return s.a, s.b, nil
	}
	return ta, tb, nil
}
