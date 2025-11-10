package injector

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
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

func (s *StdInjector) checkValidOverride(
	key types.BindKey, shouldOverride, wasOverridden bool,
) error {
	if wasOverridden && (!s.cacheOpts.Is(injopts.CacheOptAllowOverride) || !shouldOverride) {
		return remyErrs.ErrAlreadyBound{Key: key}
	}
	return nil
}

func (s *StdInjector) BindElem(bType types.BindKey, value any, opts types.BindOptions) (err error) {
	var wasOverridden bool
	if opts.Tag == "" {
		wasOverridden, err = s.cacheStorage.Set(bType, value)
	} else {
		wasOverridden, err = s.cacheStorage.SetNamed(bType, opts.Tag, value)
	}
	if err != nil {
		return err
	}

	return s.checkValidOverride(bType, opts.SoftOverride, wasOverridden)
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

	resultList = append(resultList, parentElements...)

	return
}
