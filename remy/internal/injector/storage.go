package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

// ElementsStorage holds all dependencies
type (
	genericAnyMap[T comparable]   map[T]any
	ElementsStorage[T comparable] struct {
		opts          injopts.CacheConfOption
		reflectOpts   types.ReflectionOptions
		namedElements map[string]genericAnyMap[T]
		elements      genericAnyMap[T]
	}
)

func NewElementsStorage[T comparable](
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

func (s *ElementsStorage[T]) ReflectOpts() types.ReflectionOptions {
	return s.reflectOpts
}

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool) {
	if _, ok := s.elements[key]; ok {
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
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
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
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

func (s *ElementsStorage[T]) GetAll(optKey ...string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		err = utils.ErrElementNotRegistered
		return
	}

	fromList := s.elements
	if len(optKey) > 0 {
		fromList = s.namedElements[optKey[0]]
	}

	resultList = make([]any, len(fromList))
	for _, value := range fromList {
		resultList = append(resultList, value)
	}
	return
}
