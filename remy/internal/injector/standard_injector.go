package injector

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

type (
	StdInjector struct {
		cacheOpts      injopts.CacheConfOption
		reflectOpts    types.ReflectionOptions
		parentInjector types.DependencyRetriever
		cacheStorage   types.Storage[types.BindKey]
	}
)

func New(
	opts injopts.CacheConfOption,
	reflectOpts types.ReflectionOptions,
	parent ...types.DependencyRetriever,
) *StdInjector {
	var parentInjector types.DependencyRetriever
	if len(parent) > 0 {
		parentInjector = parent[0]
	}

	return &StdInjector{
		cacheOpts:      opts,
		parentInjector: parentInjector,
		reflectOpts:    reflectOpts,
		cacheStorage:   NewElementsStorage[types.BindKey](opts, reflectOpts),
	}
}

func (s *StdInjector) SubInjector(overrides ...bool) types.Injector {
	var canOverride bool
	if len(overrides) > 0 {
		canOverride = overrides[0]
	}

	subOpts := s.cacheOpts
	if canOverride {
		subOpts |= injopts.CacheOptAllowOverride
	} else if subOpts.Is(injopts.CacheOptAllowOverride) {
		subOpts -= injopts.CacheOptAllowOverride
	}

	return New(subOpts, s.reflectOpts, s)
}

func (s *StdInjector) WrapRetriever() types.Injector {
	return nil
}

func (s *StdInjector) ReflectOpts() types.ReflectionOptions {
	return s.reflectOpts
}

func (s *StdInjector) Bind(key types.BindKey, value any) error {
	if s.cacheStorage.Set(key, value) {
		return utils.ErrAlreadyBound
	}
	return nil
}

func (s *StdInjector) BindNamed(bType types.BindKey, name string, value any) error {
	if s.cacheStorage.SetNamed(bType, name, value) {
		return utils.ErrAlreadyBound
	}
	return nil
}

func (s *StdInjector) Get(key types.BindKey) (result any, err error) {
	if result, err = s.cacheStorage.Get(key); err != nil && s.parentInjector != nil {
		cacheErr := err
		result, err = s.parentInjector.Get(key)
		if err != nil {
			err = remyErrs.ErrWrapParentSubErrors{MainError: cacheErr, SubError: err}
		}
	}
	return
}

func (s *StdInjector) GetNamed(bType types.BindKey, name string) (result any, err error) {
	if result, err = s.cacheStorage.GetNamed(bType, name); err != nil && s.parentInjector != nil {
		cacheErr := err
		result, err = s.parentInjector.GetNamed(bType, name)
		if err != nil {
			err = remyErrs.ErrWrapParentSubErrors{MainError: cacheErr, SubError: err}
		}
	}
	return
}

func (s *StdInjector) GetAll(optKey ...string) (resultList []any, err error) {
	var (
		cachedElements []any
		parentElements []any
	)

	if cachedElements, err = s.cacheStorage.GetAll(optKey...); err != nil {
		return
	}

	if s.parentInjector != nil {
		if parentElements, err = s.parentInjector.GetAll(optKey...); err != nil {
			return
		}
	}

	resultList = make([]any, len(cachedElements), len(cachedElements)+len(parentElements))
	copy(resultList, cachedElements)

	for _, element := range parentElements {
		resultList = append(resultList, element)
	}

	return
}
