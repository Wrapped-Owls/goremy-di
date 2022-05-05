package injector

import (
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
)

type StdInjector struct {
	dynamicDependencies map[types.BindKey]any
	storage             *storage.DepsStorage
}

func NewStdInjector() *StdInjector {
	return &StdInjector{
		dynamicDependencies: map[types.BindKey]any{},
		storage:             storage.NewDepsStorage(),
	}
}

func (s StdInjector) Get() any {
	//TODO implement me
	panic("implement me")
}

func (s StdInjector) GetNamed() any {
	//TODO implement me
	panic("implement me")
}
