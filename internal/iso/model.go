package iso

import "fmt"

// Model is a validated, reusable description of one single-component Langmuir
// isotherm. Construct it once with NewModel and then evaluate it at many
// driving values; the validation happens at construction time so the hot path
// (Coverage / Adsorption / Point) stays free of error handling.
type Model struct {
	K    float64
	Qmax float64
	Mode Driver
}

// NewModel validates the constants and returns a Model. A zero K or qmax, or an
// unknown driver, is rejected so that downstream callers can trust the values.
func NewModel(K, qmax float64, mode Driver) (Model, error) {
	if K <= 0 {
		return Model{}, fmt.Errorf("%w: got k=%v", ErrNonPositiveK, K)
	}
	if qmax <= 0 {
		return Model{}, fmt.Errorf("%w: got qmax=%v", ErrNonPositiveQmax, qmax)
	}
	if mode != Pressure && mode != Concentration {
		return Model{}, fmt.Errorf("unknown driver %v", mode)
	}
	return Model{K: K, Qmax: qmax, Mode: mode}, nil
}

// Coverage evaluates the coverage at driving value x. It does not re-validate
// (x < 0 is allowed to flow through the formula so callers can probe limiting
// behaviour, but for physical inputs x is non-negative).
func (m Model) Coverage(x float64) float64 {
	return Coverage(m.K, x)
}

// Adsorption evaluates the adsorbed amount at driving value x.
func (m Model) Adsorption(x float64) float64 {
	return Adsorption(m.Coverage(x), m.Qmax)
}

// HalfPressure returns x_{1/2} = 1/K for this model.
func (m Model) HalfPressure() float64 {
	return HalfPressure(m.K)
}

// Point computes a single sample at driving value x.
func (m Model) Point(x float64) Point {
	t := m.Coverage(x)
	return Point{X: x, Theta: t, Q: Adsorption(t, m.Qmax)}
}

// Curve returns the sampled isotherm described by s. It reuses Points with this
// model's constants.
func (m Model) Curve(s Scan) ([]Point, error) {
	return Points(m.K, m.Qmax, s)
}

// Report is a human-readable summary of one evaluation point. The CLI compute
// command prints it so the user sees the numbers the solver produced.
type Report struct {
	K       float64
	Qmax    float64
	X       float64
	Mode    Driver
	Theta   float64
	Q       float64
	XHalf   float64
	HenryOK bool
}

// Evaluate fills a Report for the given driving value, also recording the
// half-coverage pressure and whether the Henry low-pressure law holds to within
// a relative tolerance (used purely for display, not for validation).
func (m Model) Evaluate(x float64, henryTol float64) Report {
	r := Report{
		K:     m.K,
		Qmax:  m.Qmax,
		X:     x,
		Mode:  m.Mode,
		Theta: m.Coverage(x),
		Q:     m.Adsorption(x),
		XHalf: m.HalfPressure(),
	}
	if x > 0 {
		r.HenryOK = HenryRelativeError(m.K, x) <= henryTol
	}
	return r
}

// Summary returns a one-line description of the isotherm constants.
func (m Model) Summary() string {
	return fmt.Sprintf("Langmuir[%s]: K=%.6g qmax=%.6g x_{1/2}=%.6g",
		m.Mode, m.K, m.Qmax, m.HalfPressure())
}
