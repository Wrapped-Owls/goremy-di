package benchparse

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/perf/benchfmt"

	"github.com/wrapped-owls/goremy-di/benchmark/internal/benchname"
)

func ParseLibraries(paths []string) ([]LibraryData, []string, error) {
	libraries := make([]LibraryData, 0, len(paths))
	order := make([]string, 0)
	seenBench := map[string]struct{}{}
	labelCounts := map[string]int{}

	for _, path := range paths {
		label := uniqueLabel(baseNameNoExt(path), labelCounts)
		lib := LibraryData{
			Name:    benchname.HumanizeLibrary(label),
			Benches: map[string]Metrics{},
		}

		if err := parseLibraryFile(path, &lib, seenBench, &order); err != nil {
			return nil, nil, err
		}

		libraries = append(libraries, lib)
	}

	if len(order) == 0 {
		order = allBenchNames(libraries)
	}

	return libraries, order, nil
}

func PickMeta(libraries []LibraryData) (goos, goarch, cpu string) {
	for _, lib := range libraries {
		if goos == "" && lib.Goos != "" {
			goos = lib.Goos
		}
		if goarch == "" && lib.Goarch != "" {
			goarch = lib.Goarch
		}
		if cpu == "" && lib.CPU != "" {
			cpu = lib.CPU
		}
	}
	return
}

func parseLibraryFile(
	path string,
	lib *LibraryData,
	seenBench map[string]struct{},
	order *[]string,
) error {
	files := benchfmt.Files{Paths: []string{path}}
	for files.Scan() {
		switch rec := files.Result().(type) {
		case *benchfmt.SyntaxError:
			return fmt.Errorf("%s", rec)
		case *benchfmt.Result:
			benchName := normalizeBenchmarkName(rec.Name.String())
			if benchName == "" {
				continue
			}

			setMeta(lib, rec)
			lib.Benches[benchName] = extractMetrics(rec)
			if _, ok := seenBench[benchName]; !ok {
				seenBench[benchName] = struct{}{}
				*order = append(*order, benchName)
			}
		}
	}
	return files.Err()
}

func setMeta(lib *LibraryData, result *benchfmt.Result) {
	if lib.Goos == "" {
		lib.Goos = strings.TrimSpace(result.GetConfig("goos"))
	}
	if lib.Goarch == "" {
		lib.Goarch = strings.TrimSpace(result.GetConfig("goarch"))
	}
	if lib.CPU == "" {
		lib.CPU = strings.TrimSpace(result.GetConfig("cpu"))
	}
}

func uniqueLabel(base string, counts map[string]int) string {
	counts[base]++
	if counts[base] == 1 {
		return base
	}
	return fmt.Sprintf("%s_%d", base, counts[base])
}

func baseNameNoExt(path string) string {
	base := filepath.Base(path)
	if idx := strings.LastIndex(base, "."); idx > 0 {
		base = base[:idx]
	}
	base = strings.TrimSpace(base)
	if base == "" {
		return "unknown"
	}
	return base
}

func normalizeBenchmarkName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	if idx := strings.LastIndex(name, "-"); idx > 0 {
		if _, err := strconv.Atoi(name[idx+1:]); err == nil {
			name = name[:idx]
		}
	}

	return name
}

func extractMetrics(result *benchfmt.Result) Metrics {
	m := Metrics{Ops: int64(result.Iters), NsOp: math.NaN(), BOp: math.NaN(), Allocs: math.NaN()}

	for _, value := range result.Values {
		unit := value.Unit
		val := value.Value
		if value.OrigUnit != "" {
			unit = value.OrigUnit
			val = value.OrigValue
		}

		switch unit {
		case "ns/op":
			m.NsOp = val
		case "sec/op":
			m.NsOp = val * 1e9
		case "B/op":
			m.BOp = val
		case "allocs/op":
			m.Allocs = val
		}
	}

	return m
}

func allBenchNames(libraries []LibraryData) []string {
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
