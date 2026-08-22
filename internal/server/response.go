package server

import "math"

// round6 trims a float to six decimal places for compact, stable JSON output.
// The solver output is well-conditioned so this never loses physical
// significance; it only keeps the wire format readable.
func round6(v float64) float64 {
	const p = 1e6
	return math.Round(v*p) / p
}

// isoResponse is the payload of POST /api/iso.
type isoResponse struct {
	K     float64 `json:"k"`
	Qmax  float64 `json:"qmax"`
	X     float64 `json:"x"`
	Mode  string  `json:"mode"`
	Theta float64 `json:"theta"`
	Q     float64 `json:"q"`
	PHalf float64 `json:"phalf"`
}

// curvePoint is one sample in a POST /api/curve response.
type curvePoint struct {
	X     float64 `json:"x"`
	Theta float64 `json:"theta"`
	Q     float64 `json:"q"`
}

// curveResponse is the payload of POST /api/curve. The Points slice is what the
// front-end plots, so it always comes from the solver, never from a static
// curve definition.
type curveResponse struct {
	K      float64      `json:"k"`
	Qmax   float64      `json:"qmax"`
	Scale  string       `json:"scale"`
	Mode   string       `json:"mode"`
	Points []curvePoint `json:"points"`
}

// compResponse is the payload of POST /api/comp.
type compResponse struct {
	ThetaA float64 `json:"theta_a"`
	ThetaB float64 `json:"theta_b"`
	Total  float64 `json:"total"`
	QA     float64 `json:"q_a"`
	QB     float64 `json:"q_b"`
	QT     float64 `json:"q_total"`
}

// inverseResponse is the payload of POST /api/inverse.
type inverseResponse struct {
	K     float64 `json:"k"`
	Theta float64 `json:"theta"`
	X     float64 `json:"x"`
}

// exampleResponse wraps the bundled example dataset as served by GET
// /api/example.
type exampleResponse struct {
	Name         string    `json:"name"`
	Adsorbate    string    `json:"adsorbate"`
	TemperatureK float64   `json:"temperature_k"`
	K            float64   `json:"k"`
	Qmax         float64   `json:"qmax"`
	Mode         string    `json:"mode"`
	Unit         string    `json:"unit"`
	Pressures    []float64 `json:"pressures"`
}
