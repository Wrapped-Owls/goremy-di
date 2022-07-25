package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

type (
	StdInjector struct {
		allowOverride   bool
		reflectOpts     types.ReflectionOptions
		parentInjector  types.DependencyRetriever
		bindStorage     types.Storage[types.BindKey]
		instanceStorage types.Storage[types.BindKey]
	}
)

func New(canOverride bool, reflectOpts types.ReflectionOptions, parent ...types.Injector) *StdInjector {
	var parentInjector types.Injector
	if len(parent) > 0 {
		parentInjector = parent[0]
	}
	return &StdInjector{
		allowOverride:   canOverride,
		parentInjector:  parentInjector,
		reflectOpts:     reflectOpts,
		bindStorage:     NewElementsStorage[types.BindKey](canOverride, reflectOpts),
		instanceStorage: NewElementsStorage[types.BindKey](canOverride, reflectOpts),
	}
}

func (s *StdInjector) SubInjector(overrides ...bool) types.Injector {
	canOverride := s.allowOverride
	if len(overrides) > 0 {
		canOverride = overrides[0]
	}

	return New(canOverride, s.reflectOpts, s)
}

func (s StdInjector) ReflectOpts() types.ReflectionOptions {
	return s.reflectOpts
}

func (s *StdInjector) Bind(key types.BindKey, value any) {
	s.bindStorage.Set(key, value)
}

func (s *StdInjector) BindNamed(name string, bType types.BindKey, value any) {
	s.bindStorage.SetNamed(bType, name, value)
}

func (s StdInjector) RetrieveBind(key types.BindKey) (result any, err error) {
	if result, err = s.bindStorage.Get(key); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.RetrieveBind(key)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}

func (s StdInjector) RetrieveNamedBind(name string, bType types.BindKey) (result any, err error) {
	if result, err = s.bindStorage.GetNamed(bType, name); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.RetrieveNamedBind(bType, name)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}

func (s *StdInjector) Set(key types.BindKey, value any) {
	s.instanceStorage.Set(key, value)
}

func (s *StdInjector) SetNamed(elementType types.BindKey, name string, value any) {
	s.instanceStorage.SetNamed(elementType, name, value)
}

func (s StdInjector) GetNamed(bindKey types.BindKey, name string) (result any, err error) {
	if result, err = s.instanceStorage.GetNamed(bindKey, name); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.GetNamed(bindKey, name)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}

func (s StdInjector) Get(bindKey types.BindKey) (result any, err error) {
	if result, err = s.instanceStorage.Get(bindKey); err != nil && s.parentInjector != nil {
		result, err = s.parentInjector.Get(bindKey)
		if err != nil {
			err = utils.ErrNoElementFoundInsideOrParent
		}
	}
	return
}
