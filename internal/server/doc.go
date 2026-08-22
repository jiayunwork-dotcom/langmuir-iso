// Package server exposes the Langmuir adsorption kernels over HTTP. It is a thin
// adapter: it parses JSON requests, delegates the arithmetic to the iso, comp
// and inverse packages, and renders JSON responses (or JSON error bodies) so
// that the web front-end and any external client always see numbers produced by
// the Go solver rather than anything hard-coded in the browser.
//
// Routes:
//
//	POST /api/iso      single-component coverage and amount at one pressure
//	POST /api/curve     sampled isotherm (points the front-end plots)
//	POST /api/comp      binary competitive coverages at one state
//	POST /api/inverse   driving value that yields a target coverage
//	GET  /api/example   the bundled example dataset
//	/                   static web UI served from the web/ directory
//
// Every error path returns a JSON object {"error": "...", "code": "..."} with a
// non-2xx status code, so a failed query is visible in the body instead of being
// swallowed into a 200 with empty data.
package server
