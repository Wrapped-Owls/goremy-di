package injector

import (
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
)

type (
	bindsMap    = map[types.BindKey]any
	StdInjector struct {
		dynamicDependencies bindsMap
		namedDynamic        map[string]bindsMap
		storage             *storage.DepsStorage

		// TODO: Add allowOverride

		// TODO: Create method to init child injectors with same access
	}
)

func New() *StdInjector {
	return &StdInjector{
		dynamicDependencies: bindsMap{},
		namedDynamic:        map[string]bindsMap{},
		storage:             storage.NewDepsStorage(),
	}
}

func (s StdInjector) Storage() types.Storage {
	return s.storage
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	s.dynamicDependencies[key] = value
}

func (s *StdInjector) BindNamed(name string, bType types.BindKey, value any) {
	var namedBinds bindsMap
	if elementMap, ok := s.namedDynamic[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = bindsMap{}
	}

	namedBinds[bType] = value
	s.namedDynamic[name] = namedBinds
}

func (s StdInjector) RetrieveBind(key types.BindKey) (any, bool) {
	result, ok := s.dynamicDependencies[key]
	return result, ok
}

func (s StdInjector) RetrieveNamedBind(name string, bType types.BindKey) (any, bool) {
	if elementMap, ok := s.namedDynamic[name]; ok && elementMap != nil {
		result, subOk := elementMap[bType]
		return result, subOk
	}
	return nil, false
}

func (s StdInjector) Get(key string) any {
	return s.Storage().Get(key)
}

func (s StdInjector) Binds() []any {
	return s.Storage().Binds()
}
