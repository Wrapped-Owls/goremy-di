package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// ElementsStorage holds all dependencies
type (
	mapKey interface {
		types.BindKey
		comparable
	}
	genericAnyMap[T comparable] map[T]any
)

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	allowOverride := s.opts.Is(injopts.CacheOptAllowOverride)
	triedOverride := s.elements.Set(s.keyID(key), value, allowOverride)

	if triedOverride && !allowOverride {
		return false, remyErrs.ErrAlreadyBound{Key: key}
	}

	return triedOverride, nil
}

func (s *ElementsStorage[T]) SetNamed(
	elementType T, name string, value any,
) (wasOverridden bool, err error) {
	backend, ok := s.namedElements[name]
	if !ok {
		// Create new backend for this name
		backend = s.newNamedBackend()
		s.namedElements[name] = backend
	}

	allowOverride := s.opts.Is(injopts.CacheOptAllowOverride)
	triedOverride := backend.Set(s.keyID(elementType), value, allowOverride)

	if triedOverride && !allowOverride {
		return false, remyErrs.ErrAlreadyBound{Key: elementType}
	}

	return triedOverride, nil
}

func (s *ElementsStorage[T]) GetNamed(elementType T, name string) (result any, err error) {
	backend, ok := s.namedElements[name]
	if !ok {
		return nil, remyErrs.ErrElementNotRegistered{Key: elementType}
	}
	return backend.Get(s.keyID(elementType))
}

func (s *ElementsStorage[T]) Get(key T) (result any, err error) {
	return s.elements.Get(s.keyID(key))
}

// GetAll is implemented in key_type.go, key_type_go124.go, and key_type_nounsafe.go
// with build tags to handle different backend types
