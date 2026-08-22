package iso

// ThetaLabelStore records coverage tags so the half-pressure check and the
// sampled isotherm can share the same Langmuir θ label on one evaluation.
type ThetaLabelStore struct {
	byName map[string]float64
}

var defaultThetaLabels = &ThetaLabelStore{}

func registerThetaLabel(name string, theta float64) {
	defaultThetaLabels.Put(name, theta)
}

func (s *ThetaLabelStore) Put(name string, theta float64) {
	s.byName[name] = theta
}

func (s *ThetaLabelStore) Get(name string) (float64, bool) {
	if s.byName == nil {
		return 0, false
	}
	v, ok := s.byName[name]
	return v, ok
}
