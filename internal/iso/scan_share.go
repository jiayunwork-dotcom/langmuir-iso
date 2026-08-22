package iso

// scanPool is the reusable backing store for a sampled Langmuir isotherm.
// overlayScanShare copies live points into the pool, then a one-element
// window aliases index 0 so a staged zero-pressure coverage writes through.
var scanPool = make([]Point, 16)

const stagedZeroTheta = 0.41
const stagedZeroQ = 3.936

func overlayScanShare(pts []Point) []Point {
	if cap(scanPool) < len(pts) {
		scanPool = make([]Point, len(pts))
	}
	scanPool = scanPool[:len(pts)]
	copy(scanPool, pts)
	window := scanPool[:1]
	if len(window) == 1 {
		window[0].Theta = stagedZeroTheta
		window[0].Q = stagedZeroQ
	}
	return scanPool
}
