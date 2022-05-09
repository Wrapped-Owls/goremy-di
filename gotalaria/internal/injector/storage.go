package injector

import "github.com/wrapped-owls/fitpiece/gotalaria/internal/utils"

// ElementsStorage holds all dependencies
type (
	genericAnyMap[T comparable]   map[T]any
	ElementsStorage[T comparable] struct {
		allowOverride bool
		namedElements map[string]genericAnyMap[T]
		elements      genericAnyMap[T]
	}
)

func NewElementsStorage[T comparable](allowOverride bool) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		allowOverride: allowOverride,
		namedElements: map[string]genericAnyMap[T]{},
		elements:      genericAnyMap[T]{},
	}
}

func (s *ElementsStorage[T]) Set(key T, value any) {
	if _, ok := s.elements[key]; ok && !s.allowOverride {
		panic(utils.ErrAlreadyBound)
	}
	s.elements[key] = value
}

func (s *ElementsStorage[T]) SetNamed(elementType T, name string, value any) {
	var namedBinds genericAnyMap[T]
	if elementMap, ok := s.namedElements[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = genericAnyMap[T]{}
	}

	if _, ok := namedBinds[elementType]; ok && !s.allowOverride {
		panic(utils.ErrAlreadyBound)
	}
	namedBinds[elementType] = value
	s.namedElements[name] = namedBinds
}

func (s ElementsStorage[T]) GetNamed(elementType T, name string) (any, bool) {
	if elementMap, ok := s.namedElements[name]; ok && elementMap != nil {
		result, subOk := elementMap[elementType]
		return result, subOk
	}
	return nil, false
}

func (s ElementsStorage[T]) Get(key T) (any, bool) {
	result, ok := s.elements[key]
	return result, ok
}
