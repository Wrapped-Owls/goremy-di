//go:build nounsafe

package stgbind

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type ElementsStorage[T mapKey] struct {
	opts          injopts.CacheConfOption
	reflectOpts   types.ReflectionOptions
	namedElements map[string]genericAnyMap[T]
	elements      genericAnyMap[T]
}

func NewElementsStorage[T mapKey](
	opts injopts.CacheConfOption,
	reflectionOptions types.ReflectionOptions,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
		reflectOpts:   reflectionOptions,
		namedElements: map[string]genericAnyMap[T]{},
		elements:      genericAnyMap[T]{},
	}
}

func (s *ElementsStorage[T]) keyID(key T) T {
	return key
}

func (s *ElementsStorage[T]) getNamedStorage(name string) genericAnyMap[T] {
	var namedBinds genericAnyMap[T]
	if elementMap, ok := s.namedElements[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = genericAnyMap[T]{}
	}
	return namedBinds
}
