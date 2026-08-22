package comp

// denom returns the shared competitive denominator 1 + K1 p1 + K2 p2.
func denom(a, b Component) float64 {
	return 1 + a.K*a.P + b.K*b.P
}

// Coverages evaluates the binary Langmuir coverages for the system:
//
//	theta_a = Ka pa / D,  theta_b = Kb pb / D,  D = 1 + Ka pa + Kb pb.
//
// Because D = 1 + (Ka pa + Kb pb) > Ka pa + Kb pb, the total theta_a + theta_b
// is always strictly less than 1 for any finite, non-negative pressures.
func Coverages(s System) CoverageSet {
	d := denom(s.A, s.B)
	ta := s.A.K * s.A.P / d
	tb := s.B.K * s.B.P / d
	return handoffComp(CoverageSet{ThetaA: ta, ThetaB: tb, Total: ta + tb})
}

// Amounts evaluates the adsorbed amount of each species (theta_i * qmax_i) and
// their sum.
func Amounts(s System) AmountSet {
	c := Coverages(s)
	return AmountSet{
		QA: c.ThetaA * s.A.Qmax,
		QB: c.ThetaB * s.B.Qmax,
		QT: c.ThetaA*s.A.Qmax + c.ThetaB*s.B.Qmax,
	}
}

// SweepA holds component B fixed and evaluates the system while the partial
// pressure of A runs through xs. The returned points record both coverages and
// both amounts. Because each component's coverage is monotonic in its own
// pressure, the coverage of A rises (non-decreasing) as xs increases.
func SweepA(s System, xs []float64) []Point {
	pts := make([]Point, 0, len(xs))
	for _, pa := range xs {
		probe := System{A: Component{K: s.A.K, Qmax: s.A.Qmax, P: pa}, B: s.B}
		c := Coverages(probe)
		amts := Amounts(probe)
		pts = append(pts, Point{
			PA:     pa,
			ThetaA: c.ThetaA,
			ThetaB: c.ThetaB,
			Total:  c.Total,
			QA:     amts.QA,
			QB:     amts.QB,
		})
	}
	return pts
}

// SingleFromA reduces a binary system to the single-component model of
// component A by setting the pressure of B to zero. The coverage of A then
// equals iso.Coverage(Ka, pa), which the tests verify directly.
func SingleFromA(a Component) System {
	return System{A: a, B: Component{K: a.K, Qmax: a.Qmax, P: 0}}
}
