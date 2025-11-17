package main

type TypeInfo struct {
	Name         string
	Type         string
	Constructor  string
	Dependencies []string
	VarName      string
}

var typesInfo = map[string]TypeInfo{
	"Config": {
		Name:        "Config",
		Type:        "*fixtures.Config",
		Constructor: "fixtures.NewConfig",
		VarName:     "cfg",
	},
	"Repository": {
		Name:         "Repository",
		Type:         "fixtures.Repository",
		Constructor:  "fixtures.NewRepository",
		Dependencies: []string{"Config"},
		VarName:      "repo",
	},
	"ServiceA": {
		Name:         "ServiceA",
		Type:         "fixtures.ServiceA",
		Constructor:  "fixtures.NewServiceA",
		Dependencies: []string{"Repository", "Config"},
		VarName:      "serviceA",
	},
	"ServiceB": {
		Name:         "ServiceB",
		Type:         "fixtures.ServiceB",
		Constructor:  "fixtures.NewServiceB",
		Dependencies: []string{"ServiceA", "Config"},
		VarName:      "serviceB",
	},
	"ServiceC": {
		Name:         "ServiceC",
		Type:         "fixtures.ServiceC",
		Constructor:  "fixtures.NewServiceC",
		Dependencies: []string{"ServiceA", "ServiceB", "Config"},
		VarName:      "serviceC",
	},
}
