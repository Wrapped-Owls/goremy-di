package main

type TypeInfo struct {
	Name         string
	Type         string
	Constructor  string
	VarName      string
	Dependencies []string
}

const (
	varTypeString      = "String"
	varTypeStructEmpty = "EmptyStruct"
	varTypeBoolean     = "Boolean"
	varTypeByteArray   = "ByteArray"
	varTypeConfig      = "Config"
	varTypeRepository  = "Repository"
	varTypeServiceA    = "ServiceA"
	varTypeServiceB    = "ServiceB"
	varTypeServiceC    = "ServiceC"
)

var typesInfo = map[string]TypeInfo{
	varTypeString: {
		Name:        varTypeString,
		Type:        "string",
		Constructor: `func() string { return "string?" }`,
		VarName:     "sampleSTR",
	},
	varTypeBoolean: {
		Name:        varTypeBoolean,
		Type:        "bool",
		Constructor: "func() bool { return true }",
		VarName:     "sampleBOOL",
	},
	varTypeStructEmpty: {
		Name:        varTypeStructEmpty,
		Type:        "struct{}",
		Constructor: "func() (empty struct{}) { return empty }",
		VarName:     "sampleEmptyStruct",
	},
	varTypeByteArray: {
		Name:        varTypeByteArray,
		Type:        "[]byte",
		Constructor: "func() []byte { return []byte{} }",
		VarName:     "sampleByteArr",
	},
	varTypeConfig: {
		Name:        varTypeConfig,
		Type:        "*fixtures.Config",
		Constructor: "fixtures.NewConfig",
		VarName:     "cfg",
	},
	varTypeRepository: {
		Name:         varTypeRepository,
		Type:         "fixtures.Repository",
		Constructor:  "fixtures.NewRepository",
		Dependencies: []string{varTypeConfig},
		VarName:      "repo",
	},
	varTypeServiceA: {
		Name:         varTypeServiceA,
		Type:         "fixtures.ServiceA",
		Constructor:  "fixtures.NewServiceA",
		Dependencies: []string{varTypeRepository, varTypeConfig},
		VarName:      "serviceA",
	},
	varTypeServiceB: {
		Name:         varTypeServiceB,
		Type:         "fixtures.ServiceB",
		Constructor:  "fixtures.NewServiceB",
		Dependencies: []string{varTypeServiceA, varTypeConfig},
		VarName:      "serviceB",
	},
	varTypeServiceC: {
		Name:         varTypeServiceC,
		Type:         "fixtures.ServiceC",
		Constructor:  "fixtures.NewServiceC",
		Dependencies: []string{varTypeServiceA, varTypeServiceB, varTypeConfig},
		VarName:      "serviceC",
	},
}
