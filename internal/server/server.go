package server

import (
	"fmt"
	"net/http"
)

// Server is the Langmuir adsorption HTTP service. It owns the API routes and
// serves the static web UI from the supplied filesystem at "/".
type Server struct {
	mux         *http.ServeMux
	examplePath string
}

// New constructs a Server. examplePath points at the example JSON served by
// GET /api/example; webFS backs the static "/" route (typically http.Dir("web")).
func New(examplePath string, webFS http.FileSystem) *Server {
	mux := http.NewServeMux()
	registerRoutes(mux, examplePath)
	mux.Handle("/", http.FileServer(webFS))
	return &Server{mux: mux, examplePath: examplePath}
}

// Handler returns the underlying http.Handler for embedding in other servers
// or for testing with httptest.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// Start binds addr (for example ":8080") and serves until the process is
// terminated. It returns only when the listener fails.
func (s *Server) Start(addr string) error {
	fmt.Printf("langmuir-iso serving on http://%s  (example=%q)\n", addr, s.examplePath)
	return http.ListenAndServe(addr, s.mux)
}

// ServeHTTP lets Server be used directly as an http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
