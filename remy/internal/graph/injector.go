package graph

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

const inlinePathCapacity = 6

// Injector tracks the resolution path: a repeat inside the current path is a cycle,
// one across siblings a diamond.
type Injector struct {
	types.Injector

	buf  [inlinePathCapacity]types.GraphNode
	path []types.GraphNode
}

func NewDetector(base types.Injector) *Injector {
	return &Injector{Injector: base}
}

// SubInjector decorates the child scope, so the chain that led here keeps tracking.
func (g *Injector) SubInjector(allowOverrides ...bool) types.Injector {
	return g.decorate(g.Injector.SubInjector(allowOverrides...))
}

// WrapScopeFactory makes a temporary scope carry the path walked so far, or a cycle
// crossing into it recurses forever instead of being reported.
func (g *Injector) WrapScopeFactory(base types.ScopeFactory) types.ScopeFactory {
	return func(
		parent types.DependencyRetriever, storage types.Storage[types.BindKey],
	) types.Injector {
		return g.decorate(base(parent, storage))
	}
}

func (g *Injector) decorate(inj types.Injector) *Injector {
	scoped := &Injector{Injector: inj}
	scoped.path = append(scoped.buf[:0], g.path...)

	return scoped
}

func (g *Injector) RetrieverFor(key types.BindKey, tag string) types.Injector {
	node := types.GraphNode{Key: key, Tag: tag}
	for _, seen := range g.path {
		if seen == node {
			cycle := make([]types.GraphNode, len(g.path), len(g.path)+1)
			copy(cycle, g.path)
			panic(&remyErrs.ErrCycleDependencyDetected{Path: append(cycle, node)})
		}
	}

	fork := &Injector{Injector: g.Injector}
	fork.path = append(append(fork.buf[:0], g.path...), node)
	return fork
}

func (g *Injector) ResolveOptions() injopts.ResolveConfOption {
	opts, _ := types.ResolveOptionsOf(g.Injector)
	return opts
}
