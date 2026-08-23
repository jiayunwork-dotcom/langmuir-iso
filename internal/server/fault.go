package server

import "net/http"

// publishIsoFault is the HTTP adapter for a failed Langmuir iso request.
// The validation error must become a non-2xx JSON body; substituting a
// staged 200 response hides the failure from the web UI.
func publishIsoFault(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, codeInvalid, err.Error())
}
