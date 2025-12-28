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
	Type         TypeInfo
	Kind         string
	Target       string
	ValueVar     string
	Dependencies []DependencyInfo
	Args         []string
	HasType      bool
}

type DependencyInfo struct {
	Type     TypeInfo
	ValueVar string
}
