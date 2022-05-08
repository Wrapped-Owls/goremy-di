package injector

import (
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
	"gotalaria/internal/utils"
)

type (
	bindsMap    = map[types.BindKey]any
	StdInjector struct {
		dynamicDependencies bindsMap
		namedDynamic        map[string]bindsMap
		instanceStorage     *storage.DepsStorage
		allowOverride       bool

		// TODO: Create method to init child injectors with same access
	}
)

func New(canOverride bool) *StdInjector {
	return &StdInjector{
		allowOverride:       canOverride,
		dynamicDependencies: bindsMap{},
		namedDynamic:        map[string]bindsMap{},
		instanceStorage:     storage.NewDepsStorage(canOverride),
	}
}

func (s StdInjector) Storage() types.Storage {
	return s.instanceStorage
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	if _, ok := s.dynamicDependencies[key]; ok && !s.allowOverride {
		panic(utils.ErrAlreadyBound)
	}
	s.dynamicDependencies[key] = value
}

func (s *StdInjector) BindNamed(name string, bType types.BindKey, value any) {
	var namedBinds bindsMap
	if elementMap, ok := s.namedDynamic[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = bindsMap{}
	}

	if _, ok := namedBinds[bType]; ok && !s.allowOverride {
		panic(utils.ErrAlreadyBound)
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
