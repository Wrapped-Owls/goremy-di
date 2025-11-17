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
			{Kind: opRegisterSingleton, Target: "Config"},
			{Kind: opRegisterSingleton, Target: "Repository"},
			{Kind: opRegisterSingleton, Target: "ServiceA"},
			{Kind: opRegisterSingleton, Target: "ServiceB"},
			{Kind: opRegisterSingleton, Target: "ServiceC"},
		},
	},
	{
		Name: "SingletonRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: "Config"},
			{Kind: opRegisterSingleton, Target: "Repository"},
			{Kind: opRegisterSingleton, Target: "ServiceA"},
			{Kind: opRegisterSingleton, Target: "ServiceB"},
			{Kind: opRegisterSingleton, Target: "ServiceC"},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: "ServiceC"},
		},
	},
	{
		Name: "FactoryRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: "Config"},
			{Kind: opRegisterFactory, Target: "Repository"},
			{Kind: opRegisterFactory, Target: "ServiceA"},
			{Kind: opRegisterFactory, Target: "ServiceB"},
			{Kind: opRegisterFactory, Target: "ServiceC"},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: "ServiceC"},
		},
	},
	{
		Name: "InstanceRetrieval",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opCreateValue, Target: "Config"},
			{Kind: opCreateValue, Target: "Repository"},
			{Kind: opCreateValue, Target: "ServiceA"},
			{Kind: opCreateValue, Target: "ServiceB"},
			{Kind: opCreateValue, Target: "ServiceC"},
			{Kind: opRegisterInstance, Target: "Config"},
			{Kind: opRegisterInstance, Target: "Repository"},
			{Kind: opRegisterInstance, Target: "ServiceA"},
			{Kind: opRegisterInstance, Target: "ServiceB"},
			{Kind: opRegisterInstance, Target: "ServiceC"},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: "ServiceC"},
		},
	},
	{
		Name: "NestedDependencyResolution",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: "Config"},
			{Kind: opRegisterSingleton, Target: "Repository"},
			{Kind: opRegisterSingleton, Target: "ServiceA"},
			{Kind: opRegisterSingleton, Target: "ServiceB"},
			{Kind: opRegisterSingleton, Target: "ServiceC"},
		},
		LoopOps: []operation{
			{Kind: opRetrieve, Target: "Config"},
			{Kind: opRetrieve, Target: "Repository"},
			{Kind: opRetrieve, Target: "ServiceA"},
			{Kind: opRetrieve, Target: "ServiceB"},
			{Kind: opRetrieve, Target: "ServiceC"},
		},
	},
	{
		Name: "MultipleRetrievals",
		SetupOps: []operation{
			{Kind: opInitContainer},
			{Kind: opRegisterSingleton, Target: "Config"},
			{Kind: opRegisterSingleton, Target: "Repository"},
			{Kind: opRegisterSingleton, Target: "ServiceA"},
			{Kind: opRegisterSingleton, Target: "ServiceB"},
			{Kind: opRegisterSingleton, Target: "ServiceC"},
		},
		LoopOps: []operation{
			{Kind: opRetrieveMultiple, Target: "ServiceC", Count: 5},
		},
	},
}
