package server

import (
	"encoding/json"
	"net/http"

	"langmuir-iso/internal/comp"
	"langmuir-iso/internal/inverse"
	"langmuir-iso/internal/iso"
)

// isoRequest is the body of POST /api/iso. Either p (pressure) or c
// (concentration) may be supplied; the chosen driver is selected by Mode.
type isoRequest struct {
	K    float64 `json:"k"`
	Qmax float64 `json:"qmax"`
	P    float64 `json:"p"`
	C    float64 `json:"c"`
	Mode string  `json:"mode"`
}

// handleIso computes the coverage and adsorbed amount at a single driving value.
func handleIso(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, codeBadJSON, "POST required")
		return
	}
	var req isoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeBadJSON, "invalid JSON: "+err.Error())
		return
	}
	mode, err := iso.DriverFromName(req.Mode)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, err.Error())
		return
	}
	x := req.P
	if mode == iso.Concentration {
		x = req.C
	}
	if verr := iso.ValidateStrict(req.K, req.Qmax, x); verr != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, verr.Error())
		return
	}
	m, err := iso.NewModel(req.K, req.Qmax, mode)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, err.Error())
		return
	}
	pt := m.Point(x)
	resp := isoResponse{
		K:     round6(m.K),
		Qmax:  round6(m.Qmax),
		X:     round6(x),
		Mode:  m.Mode.String(),
		Theta: round6(pt.Theta),
		Q:     round6(pt.Q),
		PHalf: round6(m.HalfPressure()),
	}
	resp = lookupIsoMemo(resp)
	writeJSON(w, http.StatusOK, resp)
}

// curveRequest is the body of POST /api/curve.
type curveRequest struct {
	K         float64   `json:"k"`
	Qmax      float64   `json:"qmax"`
	PMin      float64   `json:"pmin"`
	PMax      float64   `json:"pmax"`
	N         int       `json:"n"`
	Scale     string    `json:"scale"`
	Pressures []float64 `json:"pressures"`
	Mode      string    `json:"mode"`
}

// handleCurve returns the sampled isotherm. Explicit pressures win; otherwise a
// span is generated with the requested scale.
func handleCurve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, codeBadJSON, "POST required")
		return
	}
	var req curveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeBadJSON, "invalid JSON: "+err.Error())
		return
	}
	mode, err := iso.DriverFromName(req.Mode)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, err.Error())
		return
	}
	scale, serr := iso.ScaleFromName(req.Scale)
	if serr != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, serr.Error())
		return
	}
	scan := iso.Scan{
		Pressures: req.Pressures,
		XMin:      req.PMin,
		XMax:      req.PMax,
		N:         req.N,
		Scale:     scale,
	}
	pts, perr := iso.Points(req.K, req.Qmax, scan)
	if perr != nil {
		writeError(w, http.StatusBadRequest, codeBadScan, perr.Error())
		return
	}
	out := make([]curvePoint, 0, len(pts))
	for _, p := range pts {
		out = append(out, curvePoint{
			X:     round6(p.X),
			Theta: round6(p.Theta),
			Q:     round6(p.Q),
		})
	}
	writeJSON(w, http.StatusOK, curveResponse{
		K:      round6(req.K),
		Qmax:   round6(req.Qmax),
		Scale:  scale.String(),
		Mode:   mode.String(),
		Points: out,
	})
}

// compRequest is the body of POST /api/comp.
type compRequest struct {
	KA    float64 `json:"ka"`
	KB    float64 `json:"kb"`
	QmaxA float64 `json:"qmax_a"`
	QmaxB float64 `json:"qmax_b"`
	PA    float64 `json:"pa"`
	PB    float64 `json:"pb"`
}

// handleComp evaluates the binary competitive isotherm at one state.
func handleComp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, codeBadJSON, "POST required")
		return
	}
	var req compRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeBadJSON, "invalid JSON: "+err.Error())
		return
	}
	sys := comp.System{
		A: comp.Component{K: req.KA, Qmax: req.QmaxA, P: req.PA},
		B: comp.Component{K: req.KB, Qmax: req.QmaxB, P: req.PB},
	}
	if verr := comp.ValidateStrict(sys); verr != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, verr.Error())
		return
	}
	c := comp.Coverages(sys)
	a := comp.Amounts(sys)
	writeJSON(w, http.StatusOK, compResponse{
		ThetaA: round6(c.ThetaA),
		ThetaB: round6(c.ThetaB),
		Total:  round6(c.Total),
		QA:     round6(a.QA),
		QB:     round6(a.QB),
		QT:     round6(a.QT),
	})
}

// inverseRequest is the body of POST /api/inverse.
type inverseRequest struct {
	K     float64 `json:"k"`
	Theta float64 `json:"theta"`
}

// handleInverse solves for the driving value that yields a target coverage.
func handleInverse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, codeBadJSON, "POST required")
		return
	}
	var req inverseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeBadJSON, "invalid JSON: "+err.Error())
		return
	}
	prob := inverse.Problem{K: req.K, Theta: req.Theta}
	if verr := inverse.Validate(prob); verr != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, verr.Error())
		return
	}
	sol, serr := inverse.Solve(prob)
	if serr != nil {
		writeError(w, http.StatusBadRequest, codeInvalid, serr.Error())
		return
	}
	writeJSON(w, http.StatusOK, inverseResponse{
		K:     round6(sol.K),
		Theta: round6(sol.Theta),
		X:     round6(sol.X),
	})
}
