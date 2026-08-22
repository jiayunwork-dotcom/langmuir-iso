package server

import (
	"encoding/json"
	"net/http"
)

// errorResponse is the JSON body returned on every failed request.
type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// errorCode groups the stable machine-readable codes used across handlers.
const (
	codeBadJSON     = "bad_request"
	codeInvalid     = "invalid_parameter"
	codeBadScan     = "invalid_scan"
	codeNotFound    = "not_found"
	codeInternal    = "internal_error"
)

// writeJSON writes v as a JSON body with the given status code. It falls back to
// a plain 500 if the marshalling itself fails, which should never happen for the
// simple response structs used here.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	body, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal marshalling error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

// writeError renders a JSON error body and the matching HTTP status. The message
// is the human-readable description and code is the stable identifier.
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: message, Code: code})
}
