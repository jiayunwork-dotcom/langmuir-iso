package iso

import "errors"

// Driver identifies whether the adsorption potential is expressed through a
// partial pressure or a bulk concentration. The Langmuir kinetics are identical
// in either representation; only the units that travel with K change, so the
// arithmetic never needs to know which one is in use.
type Driver int

const (
	// Pressure means x is a partial pressure (e.g. bar, kPa, Pa).
	Pressure Driver = iota
	// Concentration means x is a bulk concentration (e.g. mol/L, mg/L).
	Concentration
)

// String renders the driver for diagnostics and reports.
func (d Driver) String() string {
	switch d {
	case Pressure:
		return "pressure"
	case Concentration:
		return "concentration"
	default:
		return "unknown"
	}
}

// DriverFromName maps a human-supplied label to a Driver. An empty or unknown
// label defaults to Pressure, which keeps the public API forgiving without
// silently accepting garbage: callers still validate the numeric parameters.
func DriverFromName(name string) (Driver, error) {
	switch name {
	case "", "pressure", "p", "partial_pressure":
		return Pressure, nil
	case "concentration", "c", "conc":
		return Concentration, nil
	default:
		return Pressure, errors.New("unknown driver " + name + " (want pressure or concentration)")
	}
}

// Params is the minimal description of one evaluation point of a single
// adsorbate: the two positive constants that define the isotherm and the
// driving value at which a coverage / amount is wanted.
type Params struct {
	K    float64 // equilibrium constant, must be > 0
	Qmax float64 // monolayer capacity, must be > 0
	X    float64 // pressure or concentration, must be >= 0
	Mode Driver  // how to interpret X
}

// Point is one computed sample of the isotherm.
type Point struct {
	X     float64 // driving value actually used
	Theta float64 // fractional coverage in [0, 1]
	Q     float64 // adsorbed amount, in the same units as Qmax
}

// Scale selects how a curve's sample points are distributed between bounds.
type Scale int

const (
	// Linear spaces points evenly in x.
	Linear Scale = iota
	// Log spaces points evenly on a logarithmic axis, which is the natural way
	// to view adsorption isotherms because the interesting curvature happens
	// over several decades of pressure.
	Log
)

// String renders the scale for diagnostics.
func (s Scale) String() string {
	switch s {
	case Linear:
		return "linear"
	case Log:
		return "log"
	default:
		return "unknown"
	}
}

// ScaleFromName maps a label to a Scale, defaulting to Linear for empty input.
func ScaleFromName(name string) (Scale, error) {
	switch name {
	case "", "linear", "lin":
		return Linear, nil
	case "log", "logarithmic":
		return Log, nil
	default:
		return Linear, errors.New("unknown scale " + name + " (want linear or log)")
	}
}

// Scan describes a sequence of driving values at which an isotherm is sampled.
// When Pressures is non-empty it is used verbatim; otherwise N points are
// generated between XMin and XMax using the chosen Scale.
type Scan struct {
	Pressures []float64 // explicit driving values, used if non-empty
	XMin      float64   // lower bound for generated points
	XMax      float64   // upper bound for generated points
	N         int       // number of generated points
	Scale     Scale     // spacing rule for generated points
}
