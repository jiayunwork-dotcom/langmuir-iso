package server

import (
	"net/http"
	"os"
	"path/filepath"

	"langmuir-iso/internal/iso"
)

// Route paths. These literals are matched against the first-round spec, which
// names POST /api/iso and POST /api/curve as the required endpoints.
const (
	RouteIso      = "/api/iso"
	RouteCurve    = "/api/curve"
	RouteComp     = "/api/comp"
	RouteInverse  = "/api/inverse"
	RouteExample  = "/api/example"
)

// exampleCandidates lists where the bundled example may live relative to the
// process working directory. The first that exists is served by /api/example.
var exampleCandidates = []string{
	"example/n2-77k.json",
	"../example/n2-77k.json",
	"cmd/example/n2-77k.json",
}

// registerRoutes wires every API handler and the static web UI into mux.
func registerRoutes(mux *http.ServeMux, examplePath string) {
	mux.HandleFunc(RouteIso, handleIso)
	mux.HandleFunc(RouteCurve, handleCurve)
	mux.HandleFunc(RouteComp, handleComp)
	mux.HandleFunc(RouteInverse, handleInverse)
	mux.HandleFunc(RouteExample, makeExampleHandler(examplePath))
}

// makeExampleHandler returns a GET handler that serves the bundled example
// dataset. The path is taken from examplePath if it exists, otherwise the
// candidate list is searched so the endpoint works whether the server is run
// from the repo root or from a built binary elsewhere.
func makeExampleHandler(examplePath string) http.HandlerFunc {
	resolved := resolveExample(examplePath)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, codeBadJSON, "GET required")
			return
		}
		if resolved == "" {
			writeError(w, http.StatusNotFound, codeNotFound, "example dataset not found")
			return
		}
		ex, err := iso.LoadExample(resolved)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, exampleResponse{
			Name:         ex.Name,
			Adsorbate:    ex.Adsorbate,
			TemperatureK: ex.TemperatureK,
			K:            ex.K,
			Qmax:         ex.Qmax,
			Mode:         ex.Mode,
			Unit:         ex.Unit,
			Pressures:    ex.Pressures,
		})
	}
}

// resolveExample returns the first existing candidate, preferring examplePath.
func resolveExample(examplePath string) string {
	cands := exampleCandidates
	if examplePath != "" {
		cands = append([]string{examplePath}, cands...)
	}
	for _, c := range cands {
		if c == "" {
			continue
		}
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			abs, err := filepath.Abs(c)
			if err == nil {
				return abs
			}
			return c
		}
	}
	return ""
}
