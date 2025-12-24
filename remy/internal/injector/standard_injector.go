package injector

import (
	"errors"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/stgbind"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type (
	StdInjector struct {
		cacheOpts      injopts.CacheConfOption
		parentInjector types.DependencyRetriever
		cacheStorage   types.Storage[types.BindKey]
	}
)

func New(
	opts injopts.CacheConfOption,
	parent ...types.DependencyRetriever,
) *StdInjector {
	var parentInjector types.DependencyRetriever
	if len(parent) > 0 {
		parentInjector = parent[0]
	}

	return &StdInjector{
		cacheOpts:      opts,
		parentInjector: parentInjector,
		cacheStorage:   stgbind.NewElementsStorage[types.BindKey](opts),
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

	return New(subOpts, s)
}

func (s *StdInjector) WrapRetriever() types.Injector {
	return nil
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

func (s *StdInjector) RetrieveBind(bindKey types.BindKey, tag string) (result any, err error) {
	if tag == "" {
		result, err = s.cacheStorage.Get(bindKey)
	} else {
		result, err = s.cacheStorage.GetNamed(bindKey, tag)
	}

	if err != nil && s.parentInjector != nil {
		cacheErr := err
		result, err = s.parentInjector.RetrieveBind(bindKey, tag)
		if err != nil {
			err = remyErrs.ErrWrapParentSubErrors{MainError: cacheErr, SubError: err}
		}
	}
	return result, err
}

func (s *StdInjector) GetAll(optKey ...string) (resultList []any, err error) {
	var (
		cachedElements []any
		parentElements []any
	)

	if cachedElements, err = s.cacheStorage.GetAll(optKey...); err != nil &&
		// Allow not allow return all temporarily for sub-injectors
		!errors.Is(err, remyErrs.ErrConfigNotAllowReturnAll) {
		return
	}

	if s.parentInjector != nil {
		originalError := err
		if parentElements, err = s.parentInjector.GetAll(optKey...); err != nil {
			err = remyErrs.ErrWrapParentSubErrors{MainError: err}
			if originalError != nil { // Restore original error in case the parent raises an error as well
				err = originalError
			}
		}
	}
	if err != nil {
		return nil, err
	}

	resultList = make([]any, len(cachedElements), len(cachedElements)+len(parentElements))
	copy(resultList, cachedElements)

	resultList = append(resultList, parentElements...)

	return
}
