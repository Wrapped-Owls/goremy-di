package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

type CycleDetectorInjector struct {
	ij              Injector
	dependencyGraph *types.DependencyGraph
	config          Config
}

func NewCycleDetectorInjector(configs ...Config) *CycleDetectorInjector {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
	}
	return &CycleDetectorInjector{ij: NewInjector(config), config: config}
}

func (c CycleDetectorInjector) Bind(key types.BindKey, element any) error {
	return c.ij.Bind(key, element)
}

func (c CycleDetectorInjector) BindNamed(key types.BindKey, name string, element any) error {
	return c.ij.BindNamed(key, name, element)
}

func (c CycleDetectorInjector) SubInjector(allowOverrides ...bool) types.Injector {
	var shouldOverride bool
	if len(allowOverrides) > 0 {
		shouldOverride = allowOverrides[0]
	}
	inj := NewCycleDetectorInjector(Config{
		ParentInjector:     c,
		CanOverride:        shouldOverride,
		GenerifyInterfaces: c.config.GenerifyInterfaces,
		UseReflectionType:  c.config.UseReflectionType,
	})
	return inj
}

func (c CycleDetectorInjector) WrapRetriever() Injector {
	inj := NewCycleDetectorInjector(c.config)
	inj.ij = c.ij
	newGraph := types.DependencyGraph{
		UnnamedDependency: types.BindDependencies[bool]{},
		NamedDependency:   types.BindDependencies[map[string]bool]{},
	}
	if c.dependencyGraph == nil {
		inj.dependencyGraph = &newGraph
	} else {
		for key, value := range c.dependencyGraph.NamedDependency {
			nameMap := map[string]bool{}
			for name, used := range value {
				nameMap[name] = used
			}
			newGraph.NamedDependency[key] = nameMap
		}
		for key, value := range c.dependencyGraph.UnnamedDependency {
			newGraph.UnnamedDependency[key] = value
		}
		inj.dependencyGraph = &newGraph
	}
	return inj
}

func (c CycleDetectorInjector) GetNamed(key types.BindKey, name string) (any, error) {
	if c.dependencyGraph != nil {
		nameMap, ok := c.dependencyGraph.NamedDependency[key]
		if !ok {
			nameMap = map[string]bool{}
		}
		if _, hasKey := nameMap[name]; hasKey {
			panic(utils.ErrCycleDependencyDetected)
		}
		nameMap[name] = true
		c.dependencyGraph.NamedDependency[key] = nameMap
	}
	return c.ij.GetNamed(key, name)
}

func (c CycleDetectorInjector) Get(key types.BindKey) (any, error) {
	if c.dependencyGraph != nil {
		if _, hasKey := c.dependencyGraph.UnnamedDependency[key]; hasKey {
			panic(utils.ErrCycleDependencyDetected)
		} else {
			c.dependencyGraph.UnnamedDependency[key] = true
		}
	}
	return c.ij.Get(key)
}

func (c CycleDetectorInjector) ReflectOpts() types.ReflectionOptions {
	return types.ReflectionOptions{
		GenerifyInterface: c.config.GenerifyInterfaces,
		UseReflectionType: c.config.UseReflectionType,
	}
}
