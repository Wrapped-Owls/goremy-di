package injector

import (
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
)

type (
	StdInjector struct {
		bindStorage     types.Storage[types.BindKey]
		instanceStorage types.Storage[types.BindKey]
		allowOverride   bool

		// TODO: Create method to init child injectors with same access
	}
)

func New(canOverride bool) *StdInjector {
	return &StdInjector{
		allowOverride:   canOverride,
		bindStorage:     storage.NewElementsStorage[types.BindKey](canOverride),
		instanceStorage: storage.NewElementsStorage[types.BindKey](canOverride),
	}
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	s.bindStorage.Set(key, value)
}

func (s *StdInjector) BindNamed(name string, bType types.BindKey, value any) {
	s.bindStorage.SetNamed(bType, name, value)
}

func (s StdInjector) RetrieveBind(key types.BindKey) (any, bool) {
	return s.bindStorage.Get(key)
}

func (s StdInjector) RetrieveNamedBind(name string, bType types.BindKey) (any, bool) {
	return s.bindStorage.GetNamed(bType, name)
}

func (s *StdInjector) Set(key types.BindKey, value any) {
	s.instanceStorage.Set(key, value)
}

func (s *StdInjector) SetNamed(elementType types.BindKey, name string, value any) {
	s.instanceStorage.SetNamed(elementType, name, value)
}

func (s StdInjector) GetNamed(bindKey types.BindKey, name string) (any, bool) {
	return s.instanceStorage.GetNamed(bindKey, name)
}

func (s StdInjector) Get(bindKey types.BindKey) (any, bool) {
	return s.instanceStorage.Get(bindKey)
}
