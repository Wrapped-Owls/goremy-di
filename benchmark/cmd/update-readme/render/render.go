package render

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchname"
	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchparse"
)

func BuildSection(
	goVersion, goos, goarch, cpu string,
	libraries []benchparse.LibraryData,
	order []string,
) string {
	sections := make([]benchSection, 0)
	benchOrder := order
	if len(benchOrder) == 0 {
		benchOrder = allBenchNames(libraries)
	}

	for _, benchName := range benchOrder {
		best := bestForBench(libraries, benchName)
		rows := make([]benchRow, 0, len(libraries))
		for _, lib := range libraries {
			rows = append(rows, makeBenchRow(lib.Name, lib.Benches[benchName], best))
		}
		sections = append(
			sections,
			benchSection{Title: benchname.HumanizeBenchmark(benchName), Rows: rows},
		)
	}

	data := sectionTemplateData{
		MarkerStart: MarkerStart,
		MarkerEnd:   MarkerEnd,
		UpdatedUTC:  time.Now().UTC().Format(time.RFC3339),
		GoVersion:   goVersion,
		OSArch:      strings.Trim(strings.Join([]string{goos, goarch}, "/"), "/"),
		CPU:         strings.TrimSpace(cpu),
		Sections:    sections,
	}

	var out bytes.Buffer
	if err := sectionTemplate.Execute(&out, data); err != nil {
		panic(fmt.Sprintf("render benchmark section: %v", err))
	}
	return out.String()
}

func allBenchNames(libraries []benchparse.LibraryData) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0)
	for _, lib := range libraries {
		for name := range lib.Benches {
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}
