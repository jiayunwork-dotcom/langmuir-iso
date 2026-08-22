package monolayer

import "errors"

// ErrBadKinetics is returned when adsorption or desorption rate constants are
// non-positive, which removes the physical meaning of the kinetic model.
var ErrBadKinetics = errors.New("kinetic constants must be > 0")

// ErrBadPressure is returned when a negative pressure is supplied to the
// kinetic equilibrium solver.
var ErrBadPressure = errors.New("pressure must be >= 0")
