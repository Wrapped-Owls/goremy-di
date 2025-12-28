package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type (
	stgKey interface {
		types.BindKey
		comparable
	}
	genericAnyMap[T comparable] map[T]any
)

type (
	KeyValuePair[K stgKey, T any] struct {
		Key   K
		Value T
	}
	StorageBackend[K stgKey, V any] interface {
		Set(key K, value V, allowOverride bool) (triedOverride bool)
		Get(key K) (V, error)
		Size() int
		GetAll() []V
	}

	// ElementsStorage holds all dependencies
	ElementsStorage[T stgKey] struct {
		elements      StorageBackend[T, any]
		namedElements map[string]StorageBackend[T, any]
		opts          injopts.CacheConfOption
	}
)

func NewElementsStorage[T stgKey](
	opts injopts.CacheConfOption,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
		namedElements: make(map[string]StorageBackend[T, any]),
		elements:      newBackend[T](11),
	}
}

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	allowOverride := s.opts.Is(injopts.CacheOptAllowOverride)
	triedOverride := s.elements.Set(key, value, allowOverride)

	if triedOverride && !allowOverride {
		return false, remyErrs.ErrAlreadyBound{Key: key}
	}

	return triedOverride, nil
}

func (s *ElementsStorage[T]) SetNamed(
	key T, name string, value any,
) (wasOverridden bool, err error) {
	backend, ok := s.namedElements[name]
	if !ok {
		// Create new backend for this name
		backend = newBackend[T](1)
		s.namedElements[name] = backend
	}

	allowOverride := s.opts.Is(injopts.CacheOptAllowOverride)
	triedOverride := backend.Set(key, value, allowOverride)

	if triedOverride && !allowOverride {
		return false, remyErrs.ErrAlreadyBound{Key: key}
	}

	return triedOverride, nil
}

func (s *ElementsStorage[T]) GetNamed(key T, name string) (result any, err error) {
	backend, ok := s.namedElements[name]
	if !ok {
		return nil, remyErrs.ErrElementNotRegistered{Key: key}
	}
	return backend.Get(key)
}

func (s *ElementsStorage[T]) Get(key T) (result any, err error) {
	return s.elements.Get(key)
}

func (s *ElementsStorage[T]) GetAll(keyTag string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		err = remyErrs.ErrConfigNotAllowReturnAll
		return
	}

	backend := s.elements
	if keyTag != "" {
		var ok bool
		if backend, ok = s.namedElements[keyTag]; !ok {
			return nil, remyErrs.ErrElementNotRegistered{Key: keyTag}
		}
	}

	resultList = backend.GetAll()
	return
}
