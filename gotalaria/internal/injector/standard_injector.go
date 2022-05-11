package injector

import (
	"github.com/wrapped-owls/talaria-di/gotalaria/internal/types"
)

type (
	StdInjector struct {
		allowOverride   bool
		parentInjector  types.DependencyRetriever
		bindStorage     types.Storage[types.BindKey]
		instanceStorage types.Storage[types.BindKey]
	}
)

func New(canOverride bool, parent ...types.Injector) *StdInjector {
	var parentInjector types.Injector
	if len(parent) > 0 {
		parentInjector = parent[0]
	}
	return &StdInjector{
		allowOverride:   canOverride,
		parentInjector:  parentInjector,
		bindStorage:     NewElementsStorage[types.BindKey](canOverride),
		instanceStorage: NewElementsStorage[types.BindKey](canOverride),
	}
}

func (s *StdInjector) SubInjector(overrides ...bool) types.Injector {
	canOverride := s.allowOverride
	if len(overrides) > 0 {
		canOverride = overrides[0]
	}

	return New(canOverride, s)
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	s.bindStorage.Set(key, value)
}

func (s *StdInjector) BindNamed(name string, bType types.BindKey, value any) {
	s.bindStorage.SetNamed(bType, name, value)
}

func (s StdInjector) RetrieveBind(key types.BindKey) (any, bool) {
	result, ok := s.bindStorage.Get(key)
	if !ok && s.parentInjector != nil {
		return s.parentInjector.RetrieveBind(key)
	}
	return result, ok
}

func (s StdInjector) RetrieveNamedBind(name string, bType types.BindKey) (any, bool) {
	result, ok := s.bindStorage.GetNamed(bType, name)
	if !ok && s.parentInjector != nil {
		return s.parentInjector.RetrieveNamedBind(bType, name)
	}
	return result, ok
}

func (s *StdInjector) Set(key types.BindKey, value any) {
	s.instanceStorage.Set(key, value)
}

func (s *StdInjector) SetNamed(elementType types.BindKey, name string, value any) {
	s.instanceStorage.SetNamed(elementType, name, value)
}

func (s StdInjector) GetNamed(bindKey types.BindKey, name string) (any, bool) {
	result, ok := s.instanceStorage.GetNamed(bindKey, name)
	if !ok && s.parentInjector != nil {
		return s.parentInjector.GetNamed(bindKey, name)
	}
	return result, ok
}

func (s StdInjector) Get(bindKey types.BindKey) (any, bool) {
	result, ok := s.instanceStorage.Get(bindKey)
	if !ok && s.parentInjector != nil {
		return s.parentInjector.Get(bindKey)
	}
	return result, ok
}
