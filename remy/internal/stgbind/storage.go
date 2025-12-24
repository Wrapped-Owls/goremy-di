package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// ElementsStorage holds all dependencies
type (
	genericAnyMap[T comparable]   map[T]any
	ElementsStorage[T comparable] struct {
		opts          injopts.CacheConfOption
		namedElements map[string]genericAnyMap[T]
		elements      genericAnyMap[T]
	}
)

func NewElementsStorage[T comparable](
	opts injopts.CacheConfOption,
) *ElementsStorage[T] {
	return &ElementsStorage[T]{
		opts:          opts,
		namedElements: map[string]genericAnyMap[T]{},
		elements:      genericAnyMap[T]{},
	}
}

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	if _, ok := s.elements[key]; ok {
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
			return false, remyErrs.ErrAlreadyBound{Key: key}
		}
		wasOverridden = true
	}
	s.elements[key] = value
	return
}

func (s *ElementsStorage[T]) SetNamed(
	elementType T, name string, value any,
) (wasOverridden bool, err error) {
	var namedBinds genericAnyMap[T]
	if elementMap, ok := s.namedElements[name]; ok {
		namedBinds = elementMap
	} else {
		namedBinds = genericAnyMap[T]{}
	}

	if _, ok := namedBinds[elementType]; ok {
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
			return false, remyErrs.ErrAlreadyBound{Key: elementType}
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
			err = remyErrs.ErrElementNotRegistered{Key: elementType}
		}
		return result, err
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: elementType}
}

func (s *ElementsStorage[T]) Get(key T) (result any, err error) {
	var ok bool
	if result, ok = s.elements[key]; !ok {
		err = remyErrs.ErrElementNotRegistered{Key: key}
	}
	return
}

func (s *ElementsStorage[T]) GetAll(optKey ...string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		err = remyErrs.ErrConfigNotAllowReturnAll
		return
	}

	fromList := s.elements
	if len(optKey) > 0 {
		fromList = s.namedElements[optKey[0]]
	}

	resultList = make([]any, 0, len(fromList))
	for _, value := range fromList {
		resultList = append(resultList, value)
	}
	return
}
