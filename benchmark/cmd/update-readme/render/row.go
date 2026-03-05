package render

import (
	"fmt"
	"math"

	"github.com/wrapped-owls/goremy-di/benchmark/cmd/update-readme/render/benchformat"
	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchparse"
)

func makeBenchRow(name string, m benchparse.Metrics, best benchBest) benchRow {
	row := benchRow{Library: name, Ops: "-", NsOp: "-", BOp: "-", Allocs: "-"}
	if m.Ops > 0 {
		row.Ops = benchformat.Int(m.Ops)
	}
	if present(m.NsOp) {
		row.NsOp = benchformat.Number(m.NsOp)
	}
	if present(m.BOp) {
		row.BOp = benchformat.Number(m.BOp)
	}
	if present(m.Allocs) {
		row.Allocs = benchformat.Number(m.Allocs)
	}

	if m.Ops > 0 && m.Ops == best.opsMax {
		row.Ops = bold(row.Ops)
	}
	if present(m.NsOp) && nearlyEqual(m.NsOp, best.nsMin) {
		row.NsOp = bold(row.NsOp)
	}
	if present(m.BOp) && nearlyEqual(m.BOp, best.bMin) {
		row.BOp = bold(row.BOp)
	}
	if present(m.Allocs) && nearlyEqual(m.Allocs, best.allocMin) {
		row.Allocs = bold(row.Allocs)
	}

	return row
}

func bold(s string) string {
	if s == "-" {
		return s
	}
	return fmt.Sprintf("**%s**", s)
}

func nearlyEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
