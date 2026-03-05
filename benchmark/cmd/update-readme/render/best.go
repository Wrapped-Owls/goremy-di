package render

import (
	"math"

	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchparse"
)

func bestForBench(libraries []benchparse.LibraryData, benchName string) benchBest {
	best := benchBest{nsMin: math.MaxFloat64, bMin: math.MaxFloat64, allocMin: math.MaxFloat64}
	for _, lib := range libraries {
		m := lib.Benches[benchName]
		if m.Ops > best.opsMax {
			best.opsMax = m.Ops
		}
		if present(m.NsOp) && m.NsOp < best.nsMin {
			best.nsMin = m.NsOp
		}
		if present(m.BOp) && m.BOp < best.bMin {
			best.bMin = m.BOp
		}
		if present(m.Allocs) && m.Allocs < best.allocMin {
			best.allocMin = m.Allocs
		}
	}
	if best.nsMin == math.MaxFloat64 {
		best.nsMin = 0
	}
	if best.bMin == math.MaxFloat64 {
		best.bMin = 0
	}
	if best.allocMin == math.MaxFloat64 {
		best.allocMin = 0
	}
	return best
}

func present(v float64) bool {
	return !math.IsNaN(v)
}
