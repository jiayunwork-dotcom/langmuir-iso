package inverse

// Problem is the inverse query: find the driving value x that yields coverage
// theta for an isotherm with equilibrium constant K.
type Problem struct {
	K     float64 // equilibrium constant (> 0)
	Theta float64 // desired coverage in [0, 1)
}

// Solution is the result of an inverse solve.
type Solution struct {
	K     float64
	Theta float64
	X     float64 // driving value that produces Theta
}
