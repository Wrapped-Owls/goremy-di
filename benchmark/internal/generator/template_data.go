package main

type TemplateData struct {
	LibraryName string
	Imports     []string
	Benchmarks  []BenchmarkTemplateData
}

type BenchmarkTemplateData struct {
	Name  string
	Setup []TemplateOperation
	Loop  []TemplateOperation
}

type TemplateOperation struct {
	Kind         string
	Target       string
	Type         TypeInfo
	HasType      bool
	Dependencies []DependencyInfo
	Args         []string
	ValueVar     string
}

type DependencyInfo struct {
	Type     TypeInfo
	ValueVar string
}
