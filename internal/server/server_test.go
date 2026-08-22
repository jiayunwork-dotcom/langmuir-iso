package server

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAPIIsoHandler checks the /api/iso endpoint returns the correct coverage,
// adsorbed amount and half-coverage pressure for a valid request, and that the
// numbers come from the solver rather than constants.
func TestAPIIsoHandler(t *testing.T) {
	srv := New("example/n2-77k.json", http.Dir("../../web"))
	body := `{"k":0.115,"qmax":9.6,"p":5.0}`
	req := httptest.NewRequest(http.MethodPost, RouteIso, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("%s status = %d, want 200; body=%s", RouteIso, w.Code, w.Body.String())
	}
	var resp isoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	wantTheta := 0.115 * 5.0 / (1 + 0.115*5.0)
	if math.Abs(resp.Theta-wantTheta) > 1e-4 {
		t.Errorf("theta = %v, want %v", resp.Theta, wantTheta)
	}
	if math.Abs(resp.PHalf-1/0.115) > 1e-3 {
		t.Errorf("phalf = %v, want %v", resp.PHalf, 1/0.115)
	}
	if math.Abs(resp.Q-wantTheta*9.6) > 1e-3 {
		t.Errorf("q = %v, want %v", resp.Q, wantTheta*9.6)
	}
}

// TestAPICurveHandler checks the /api/curve endpoint returns one point per
// supplied pressure, with zero coverage at p=0 and a monotone-increasing
// coverage, i.e. the plotted data genuinely comes from the solver.
func TestAPICurveHandler(t *testing.T) {
	srv := New("example/n2-77k.json", http.Dir("../../web"))
	body := `{"k":0.115,"qmax":9.6,"pressures":[0,1,2,5,10],"scale":"linear"}`
	req := httptest.NewRequest(http.MethodPost, RouteCurve, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("%s status=%d; body=%s", RouteCurve, w.Code, w.Body.String())
	}
	var resp curveResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Points) != 5 {
		t.Fatalf("points = %d, want 5", len(resp.Points))
	}
	if resp.Points[0].Theta != 0 {
		t.Errorf("theta at p=0 should be 0, got %v", resp.Points[0].Theta)
	}
	for i := 1; i < len(resp.Points); i++ {
		if resp.Points[i].Theta < resp.Points[i-1].Theta {
			t.Errorf("theta not monotone at index %d", i)
		}
	}
}

// TestAPIError checks that an illegal request (K <= 0) returns a non-2xx status
// with a JSON error body, so the failure is visible to clients and the web UI
// rather than being swallowed into a 200 with empty data.
func TestAPIError(t *testing.T) {
	srv := New("example/n2-77k.json", http.Dir("../../web"))
	body := `{"k":-1,"qmax":9.6,"p":5.0}`
	req := httptest.NewRequest(http.MethodPost, RouteIso, strings.NewReader(body))
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Fatalf("%s with k<0 should not be 200; body=%s", RouteIso, w.Body.String())
	}
	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("expected error message in body, got empty")
	}
}
