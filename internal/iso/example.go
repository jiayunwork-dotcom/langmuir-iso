package iso

import (
	"encoding/json"
	"fmt"
	"os"
)

// Example is a small, self-describing dataset that documents one physical
// adsorption system and the pressures at which an isotherm should be drawn. It
// is the on-disk shape of files such as example/n2-77k.json. The file stores
// the measured system plus, optionally, pre-computed coverages for reference;
// the solver always recomputes from K and qmax so nothing is hard-coded.
type Example struct {
	Name         string    `json:"name"`
	Adsorbate    string    `json:"adsorbate"`
	TemperatureK float64   `json:"temperature_k"`
	K            float64   `json:"k"`
	Qmax         float64   `json:"qmax"`
	Mode         string    `json:"mode"`
	Unit         string    `json:"unit"`
	Pressures    []float64 `json:"pressures"`
	// Reference is an optional, pre-computed coverage per pressure kept for
	// documentation only; it is never trusted by the solver.
	Reference []float64 `json:"reference,omitempty"`
}

// LoadExample reads and validates an example file from disk. The file must
// describe a physical system (K > 0, qmax > 0) and must list at least one
// pressure.
func LoadExample(path string) (Example, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Example{}, fmt.Errorf("reading example %q: %w", path, err)
	}
	var ex Example
	if err := json.Unmarshal(raw, &ex); err != nil {
		return Example{}, fmt.Errorf("parsing example %q: %w", path, err)
	}
	if ex.K <= 0 {
		return Example{}, fmt.Errorf("%w: example %q has k=%v", ErrNonPositiveK, path, ex.K)
	}
	if ex.Qmax <= 0 {
		return Example{}, fmt.Errorf("%w: example %q has qmax=%v", ErrNonPositiveQmax, path, ex.Qmax)
	}
	if len(ex.Pressures) == 0 {
		return Example{}, fmt.Errorf("example %q lists no pressures", path)
	}
	return ex, nil
}

// Model builds a validated single-component Model from the example constants.
func (ex Example) Model() (Model, error) {
	mode, err := DriverFromName(ex.Mode)
	if err != nil {
		return Model{}, err
	}
	return NewModel(ex.K, ex.Qmax, mode)
}

// Scan builds a Scan that visits every pressure listed in the example. This is
// what the web page and the CLI compute command feed to the solver so the
// rendered curve always comes from the backend, never from the file's
// reference column.
func (ex Example) Scan() Scan {
	xs := make([]float64, len(ex.Pressures))
	copy(xs, ex.Pressures)
	return Scan{Pressures: xs}
}

// Curve solves the example isotherm at each listed pressure and returns the
// points. It is the single source of truth for any rendering of the example.
func (ex Example) Curve() ([]Point, error) {
	m, err := ex.Model()
	if err != nil {
		return nil, err
	}
	return m.Curve(ex.Scan())
}
