package injector

import (
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
)

type StdInjector struct {
	dynamicDependencies map[types.BindKey]any
	storage             *storage.DepsStorage
}

func New() *StdInjector {
	return &StdInjector{
		dynamicDependencies: map[types.BindKey]any{},
		storage:             storage.NewDepsStorage(),
	}
}

func (s StdInjector) Storage() types.Storage {
	return s.storage
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	s.dynamicDependencies[key] = value
}

func (s StdInjector) RetrieveBind(key types.BindKey) (any, bool) {
	result, ok := s.dynamicDependencies[key]
	return result, ok
}

func (s StdInjector) Get(key string) any {
	return s.Storage().Get(key)
}

func (s StdInjector) Binds() []any {
	return s.Storage().Binds()
}
