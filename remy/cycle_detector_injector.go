package remy

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// cycleDetectorInjector is the injector to be used in test file, to check if
// any of the bound dependencies have a requirement cycle
type cycleDetectorInjector struct {
	ij              Injector
	dependencyGraph *types.DependencyGraph
	config          Config
}

// NewCycleDetectorInjector creates a new Injector that is able to check for cycle dependencies during runtime.
//
// As it is much slower that the injector.StandardInjector, it is only recommended to be used in test files.
//
//goland:noinspection GoExportedFuncWithUnexportedType
func NewCycleDetectorInjector(configs ...Config) *cycleDetectorInjector {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
	}
	return &cycleDetectorInjector{ij: NewInjector(config), config: config}
}

func (c cycleDetectorInjector) BindElem(
	key types.BindKey, element any, opts types.BindOptions,
) error {
	return c.ij.BindElem(key, element, opts)
}

func (c cycleDetectorInjector) SubInjector(allowOverrides ...bool) types.Injector {
	var shouldOverride bool
	if len(allowOverrides) > 0 {
		shouldOverride = allowOverrides[0]
	}
	inj := NewCycleDetectorInjector(
		Config{
			ParentInjector: c,
			CanOverride:    shouldOverride,
		},
	)
	return inj
}

func (c cycleDetectorInjector) WrapRetriever() Injector {
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

func (c cycleDetectorInjector) RetrieveBind(bindKey types.BindKey, tag string) (any, error) {
	if c.dependencyGraph != nil {
		var hasKey bool
		if tag == "" {
			// Unnamed dependency
			_, hasKey = c.dependencyGraph.UnnamedDependency[bindKey]
			c.dependencyGraph.UnnamedDependency[bindKey] = true
		} else {
			// Named dependency
			nameMap, ok := c.dependencyGraph.NamedDependency[bindKey]
			if !ok {
				nameMap = map[string]bool{}
				c.dependencyGraph.NamedDependency[bindKey] = nameMap
			}
			_, hasKey = nameMap[tag]
			nameMap[tag] = true
		}

		if hasKey {
			panic(&remyErrs.ErrCycleDependencyDetected{Path: c.dependencyGraph})
		}
	}
	return c.ij.RetrieveBind(bindKey, tag)
}

func (c cycleDetectorInjector) GetAll(optKey ...string) ([]any, error) {
	return c.ij.GetAll(optKey...)
}
