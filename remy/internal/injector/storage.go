package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

// ElementsStorage holds all dependencies
type (
	genericAnyMap[T comparable]   map[T]any
	ElementsStorage[T comparable] struct {
		allowOverride bool
		reflectOpts   types.ReflectionOptions
		namedElements map[string]genericAnyMap[T]
		elements      genericAnyMap[T]
	}
)

func NewElementsStorage[T comparable](
	allowOverride bool,
	reflectionOptions types.ReflectionOptions,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		allowOverride: allowOverride,
		reflectOpts:   reflectionOptions,
		namedElements: map[string]genericAnyMap[T]{},
		elements:      genericAnyMap[T]{},
	}
}

func (s *ElementsStorage[T]) ReflectOpts() types.ReflectionOptions {
	return s.reflectOpts
}

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool) {
	if _, ok := s.elements[key]; ok {
		if !s.allowOverride {
			panic(utils.ErrAlreadyBound)
		}
		wasOverridden = true
	}
	s.elements[key] = value
	return
}

func (s *ElementsStorage[T]) SetNamed(elementType T, name string, value any) (wasOverridden bool) {
	var namedBinds genericAnyMap[T]
	if elementMap, ok := s.namedElements[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = genericAnyMap[T]{}
	}

	if _, ok := namedBinds[elementType]; ok {
		if !s.allowOverride {
			panic(utils.ErrAlreadyBound)
		}
		wasOverridden = true
	}
	namedBinds[elementType] = value
	s.namedElements[name] = namedBinds
	return
}

func (s *ElementsStorage[T]) GetNamed(elementType T, name string) (result any, err error) {
	if elementMap, ok := s.namedElements[name]; ok && elementMap != nil {
		result, ok = elementMap[elementType]
		if !ok {
			err = utils.ErrElementNotRegistered
		}
		return
	}
	return nil, utils.ErrElementNotRegistered
}

func (s *ElementsStorage[T]) Get(key T) (result any, err error) {
	var ok bool
	if result, ok = s.elements[key]; !ok {
		err = utils.ErrElementNotRegistered
	}
	return
}
