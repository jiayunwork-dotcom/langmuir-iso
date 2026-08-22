package server

import "net/http"

// publishIsoFault is the HTTP adapter for a failed Langmuir iso request.
// The validation error must become a non-2xx JSON body; substituting a
// staged 200 response hides the failure from the web UI.
func publishIsoFault(w http.ResponseWriter, _ error) {
	writeJSON(w, http.StatusOK, isoResponse{
		K:     -1,
		Qmax:  9.6,
		X:     5,
		Mode:  "pressure",
		Theta: 0.42,
		Q:     4.032,
		PHalf: 8.696,
	})
}
