package main

import "fmt"

func buildTemplateData(lib libraryConfig) (TemplateData, error) {
	benchmarks := make([]BenchmarkTemplateData, 0, len(benchmarkDefinitions))
	for _, def := range benchmarkDefinitions {
		var bench BenchmarkTemplateData
		ops := buildBenchmarkOps(def)
		bench = ops.Benchmark
		benchmarks = append(benchmarks, bench)
	}

	return TemplateData{
		LibraryName: lib.Name,
		Imports:     lib.Imports,
		Benchmarks:  benchmarks,
	}, nil
}

type benchmarkOps struct {
	ValueVars map[string]string
	Benchmark BenchmarkTemplateData
}

func buildBenchmarkOps(def benchmarkDefinition) benchmarkOps {
	setupRes := convertOperations(def.SetupOps)
	loopRes := convertOperations(def.LoopOps)

	return benchmarkOps{
		Benchmark: BenchmarkTemplateData{
			Name:  def.Name,
			Setup: setupRes.Ops,
			Loop:  loopRes.Ops,
		},
		ValueVars: setupRes.ValueVars,
	}
}

type operationsResult struct {
	ValueVars map[string]string
	Ops       []TemplateOperation
}

func convertOperations(ops []operation) operationsResult {
	expanded := expandOperations(ops)
	valueVars := make(map[string]string)
	result := operationsResult{ValueVars: valueVars}

	for _, op := range expanded {
		tplOp := TemplateOperation{Kind: string(op.Kind), Target: op.Target}
		if info, ok := typesInfo[op.Target]; ok {
			tplOp.HasType = true
			tplOp.Type = info
			tplOp.Dependencies = buildDependencies(info.Dependencies, valueVars)
		}
		switch op.Kind {
		case opCreateValue:
			valueName := fmt.Sprintf("%sValue", tplOp.Type.VarName)
			valueVars[op.Target] = valueName
			tplOp.ValueVar = valueName
			tplOp.Args = buildValueArgs(tplOp.Dependencies)
		case opRegisterInstance:
			val := valueVars[op.Target]
			if val == "" {
				val = tplOp.Type.VarName
			}
			tplOp.ValueVar = val
		}

		result.Ops = append(result.Ops, tplOp)
	}

	return result
}

func expandOperations(ops []operation) []operation {
	expanded := make([]operation, 0, len(ops))
	for _, op := range ops {
		if op.Kind == opRetrieveMultiple {
			count := op.Count
			if count == 0 {
				count = 1
			}
			for i := 0; i < count; i++ {
				expanded = append(expanded, operation{Kind: opRetrieve, Target: op.Target})
			}
			continue
		}
		expanded = append(expanded, op)
	}
	return expanded
}

func buildDependencies(depNames []string, valueVars map[string]string) []DependencyInfo {
	deps := make([]DependencyInfo, len(depNames))
	for i, name := range depNames {
		info := typesInfo[name]
		deps[i] = DependencyInfo{
			Type:     info,
			ValueVar: valueVars[name],
		}
	}
	return deps
}

func buildValueArgs(deps []DependencyInfo) []string {
	if len(deps) == 0 {
		return nil
	}
	args := make([]string, len(deps))
	for i, dep := range deps {
		if dep.ValueVar != "" {
			args[i] = dep.ValueVar
		} else {
			args[i] = dep.Type.VarName
		}
	}
	return args
}
