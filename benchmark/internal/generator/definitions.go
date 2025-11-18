package main

type operationKind string

const (
	opInitContainer     operationKind = "init_container"
	opRegisterSingleton operationKind = "register_singleton"
	opRegisterFactory   operationKind = "register_factory"
	opRegisterInstance  operationKind = "register_instance"
	opCreateValue       operationKind = "create_value"
	opRetrieve          operationKind = "retrieve"
	opRetrieveMultiple  operationKind = "retrieve_multiple"
)

type operation struct {
	Kind   operationKind
	Target string
	Count  int
}

type benchmarkDefinition struct {
	Name     string
	SetupOps []operation
	LoopOps  []operation
}

var benchmarkDefinitions = []benchmarkDefinition{
	{
		Name: "Registration",
		LoopOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: varTypeConfig},
			{Kind: opRegisterSingleton, Target: varTypeRepository},
			{Kind: opRegisterSingleton, Target: varTypeServiceA},
			{Kind: opRegisterSingleton, Target: varTypeServiceB},
			{Kind: opRegisterSingleton, Target: varTypeServiceC},
		},
	},
	{
		Name: "SingletonRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: varTypeConfig},
			{Kind: opRegisterSingleton, Target: varTypeRepository},
			{Kind: opRegisterSingleton, Target: varTypeServiceA},
			{Kind: opRegisterSingleton, Target: varTypeServiceB},
			{Kind: opRegisterSingleton, Target: varTypeServiceC},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: varTypeServiceC},
		},
	},
	{
		Name: "FactoryRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: varTypeConfig},
			{Kind: opRegisterFactory, Target: varTypeRepository},
			{Kind: opRegisterFactory, Target: varTypeServiceA},
			{Kind: opRegisterFactory, Target: varTypeServiceB},
			{Kind: opRegisterFactory, Target: varTypeServiceC},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: varTypeServiceC},
		},
	},
	{
		Name: "InstanceRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opCreateValue, Target: varTypeConfig},
			{Kind: opCreateValue, Target: varTypeRepository},
			{Kind: opCreateValue, Target: varTypeServiceA},
			{Kind: opCreateValue, Target: varTypeServiceB},
			{Kind: opCreateValue, Target: varTypeServiceC},
			{Kind: opRegisterInstance, Target: varTypeConfig},
			{Kind: opRegisterInstance, Target: varTypeRepository},
			{Kind: opRegisterInstance, Target: varTypeServiceA},
			{Kind: opRegisterInstance, Target: varTypeServiceB},
			{Kind: opRegisterInstance, Target: varTypeServiceC},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: varTypeServiceC},
		},
	},
	{
		Name: "NestedDependencyResolution",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: varTypeConfig},
			{Kind: opRegisterSingleton, Target: varTypeRepository},
			{Kind: opRegisterSingleton, Target: varTypeServiceA},
			{Kind: opRegisterSingleton, Target: varTypeServiceB},
			{Kind: opRegisterSingleton, Target: varTypeServiceC},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: varTypeConfig},
			{Kind: opRetrieve, Target: varTypeRepository},
			{Kind: opRetrieve, Target: varTypeServiceA},
			{Kind: opRetrieve, Target: varTypeServiceB},
			{Kind: opRetrieve, Target: varTypeServiceC},
		},
	},
	{
		Name: "MultipleRetrievals",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: varTypeConfig},
			{Kind: opRegisterSingleton, Target: varTypeRepository},
			{Kind: opRegisterSingleton, Target: varTypeServiceA},
			{Kind: opRegisterSingleton, Target: varTypeServiceB},
			{Kind: opRegisterSingleton, Target: varTypeServiceC},
		},
		LoopOps: []operation{
			{Kind: opRetrieveMultiple, Target: varTypeServiceC, Count: 5},
		},
	},
}
