//go:build nounsafe

package stgbind

import (
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

type ElementsStorage[T stgKey] struct {
	namedElements map[string]genericAnyMap[T]
	elements      genericAnyMap[T]
	opts          injopts.CacheConfOption
}

func NewElementsStorage[T stgKey](
	opts injopts.CacheConfOption,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
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
