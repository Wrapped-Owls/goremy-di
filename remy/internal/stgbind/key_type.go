//go:build !nounsafe

package stgbind

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type ElementsStorage[T mapKey] struct {
	opts          injopts.CacheConfOption
	reflectOpts   types.ReflectionOptions
	namedElements map[string]genericAnyMap[uint64]
	elements      genericAnyMap[uint64]
}

func NewElementsStorage[T mapKey](
	opts injopts.CacheConfOption,
	reflectionOptions types.ReflectionOptions,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
		reflectOpts:   reflectionOptions,
		namedElements: map[string]genericAnyMap[uint64]{},
		elements:      genericAnyMap[uint64]{},
	}
}

func (s *ElementsStorage[T]) keyID(key T) uint64 {
	return key.ID()
}

func (s *ElementsStorage[T]) getNamedStorage(name string) genericAnyMap[uint64] {
	var namedBinds genericAnyMap[uint64]
	if elementMap, ok := s.namedElements[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = genericAnyMap[uint64]{}
	}
	return namedBinds
}
