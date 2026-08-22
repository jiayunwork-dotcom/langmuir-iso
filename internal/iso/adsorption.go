package iso

// Adsorption returns the adsorbed amount q for coverage theta and monolayer
// capacity qmax:
//
//	q = theta * qmax.
//
// Because theta is in [0, 1), q is always in [0, qmax). At zero coverage the
// adsorbed amount is exactly zero; at saturation it approaches qmax from below.
func Adsorption(theta, qmax float64) float64 {
	return theta * qmax
}

// AdsorptionAt returns the adsorbed amount directly from K, qmax and the
// driving value x by composing Coverage and Adsorption in one call. It is the
// workhorse behind the /api/iso endpoint and the CLI compute command.
func AdsorptionAt(K, qmax, x float64) float64 {
	return Adsorption(Coverage(K, x), qmax)
}
