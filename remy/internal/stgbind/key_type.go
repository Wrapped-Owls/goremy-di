//go:build !nounsafe

package stgbind

import (
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type ElementsStorage[T mapKey] struct {
	namedElements map[string]genericAnyMap[uint64]
	elements      genericAnyMap[uint64]
	opts          injopts.CacheConfOption
}

func NewElementsStorage[T mapKey](
	opts injopts.CacheConfOption,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
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
