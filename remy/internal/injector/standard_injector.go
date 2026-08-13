package injector

import (
	"errors"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/stgbind"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type (
	Options struct {
		ScopeName string
		Cache     injopts.CacheConfOption
		Resolve   injopts.ResolveConfOption
	}

	StdInjector struct {
		parentInjector types.DependencyRetriever
		cacheStorage   types.Storage[types.BindKey]
		scopeName      string
		cacheOpts      injopts.CacheConfOption
		resolveOpts    injopts.ResolveConfOption
	}
)

func New(opts Options, parent ...types.DependencyRetriever) *StdInjector {
	return NewWithStorage(opts, stgbind.NewElementsStorage[types.BindKey](opts.Cache), parent...)
}

// NewWithStorage creates a StdInjector that uses the provided storage backend, so
// callers can supply an optimised one (e.g. SliceStorage for ephemeral
// sub-injectors with a known, small number of entries).
func NewWithStorage(
	opts Options,
	storage types.Storage[types.BindKey],
	parent ...types.DependencyRetriever,
) *StdInjector {
	var parentInjector types.DependencyRetriever
	if len(parent) > 0 {
		parentInjector = parent[0]
	}

	return &StdInjector{
		scopeName:      opts.ScopeName,
		cacheOpts:      opts.Cache,
		resolveOpts:    opts.Resolve,
		parentInjector: parentInjector,
		cacheStorage:   storage,
	}
}

// ScopeName returns the optional name this injector scope was created with.
func (s *StdInjector) ScopeName() string {
	return s.scopeName
}

// Parent returns the retriever lookups fall back to on misses, nil for roots.
func (s *StdInjector) Parent() types.DependencyRetriever {
	return s.parentInjector
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

	return New(Options{ScopeName: s.scopeName, Cache: subOpts, Resolve: s.resolveOpts}, s)
}

func (s *StdInjector) RetrieverFor(types.BindKey, string) types.Injector {
	return nil
}

// ResolveOptions returns this scope options merged with the inheritable ones
// coming from its ancestors.
func (s *StdInjector) ResolveOptions() injopts.ResolveConfOption {
	// what a sub-injector picks up from its ancestors; Isolated is per scope
	const inheritable = injopts.ResolveOptDuckTyping | injopts.ResolveOptTracePath

	opts := s.resolveOpts
	if s.isolated() {
		return opts
	}
	if holder, ok := s.parentInjector.(types.ResolveOptionsHolder); ok {
		opts |= holder.ResolveOptions() & inheritable
	}
	return opts
}

func (s *StdInjector) inheritsFromParent() bool {
	return s.parentInjector != nil && !s.isolated()
}

func (s *StdInjector) isolated() bool {
	return s.resolveOpts.Is(injopts.ResolveOptIsolated)
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

	if err != nil && s.inheritsFromParent() {
		cacheErr := err
		result, err = s.parentInjector.RetrieveBind(bindKey, tag)
		if err != nil {
			err = remyErrs.ErrWrapParentSubErrors{MainError: cacheErr, SubError: err}
		}
	}
	return result, err
}

func (s *StdInjector) GetAll(keyTag string) (resultList []any, err error) {
	var (
		cachedElements []any
		parentElements []any
	)

	if cachedElements, err = s.cacheStorage.GetAll(keyTag); err != nil &&
		// Allow not allow return all temporarily for sub-injectors
		!errors.Is(err, remyErrs.ErrConfigNotAllowReturnAll) {
		return
	}

	if s.inheritsFromParent() {
		originalError := err
		if parentElements, err = s.parentInjector.GetAll(keyTag); err != nil {
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
