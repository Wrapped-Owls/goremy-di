package main

const FixturesImportPath = "github.com/wrapped-owls/goremy-di/benchmark/internal/fixtures"

type libraryConfig struct {
	Name         string
	FileName     string
	TemplateName string
	Imports      []string
}

var libraries = []libraryConfig{
	{
		Name:         "Remy",
		FileName:     "remy_benchmark_test.go",
		TemplateName: "remy.go.tmpl",
		Imports: []string{
			FixturesImportPath,
			"github.com/wrapped-owls/goremy-di/remy",
		},
	},
	{
		Name:         "Do",
		FileName:     "do_benchmark_test.go",
		TemplateName: "do.go.tmpl",
		Imports: []string{
			"github.com/samber/do/v2",
			FixturesImportPath,
		},
	},
	{
		Name:         "Dig",
		FileName:     "dig_benchmark_test.go",
		TemplateName: "dig.go.tmpl",
		Imports: []string{
			FixturesImportPath,
			"go.uber.org/dig",
		},
	},
}
