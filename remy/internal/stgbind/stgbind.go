package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// ElementsStorage holds all dependencies
type (
	stgKey interface {
		types.BindKey
		comparable
	}
	genericAnyMap[T comparable] map[T]any
)

func (s *ElementsStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	if _, ok := s.elements[s.keyID(key)]; ok {
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
			return false, remyErrs.ErrAlreadyBound{Key: key}
		}
		wasOverridden = true
	}
	s.elements[s.keyID(key)] = value
	return
}

func (s *ElementsStorage[T]) SetNamed(
	elementType T, name string, value any,
) (wasOverridden bool, err error) {
	namedBinds := s.getNamedStorage(name)

	if _, ok := namedBinds[s.keyID(elementType)]; ok {
		if !s.opts.Is(injopts.CacheOptAllowOverride) {
			return false, remyErrs.ErrAlreadyBound{Key: elementType}
		}
		wasOverridden = true
	}
	namedBinds[s.keyID(elementType)] = value
	s.namedElements[name] = namedBinds
	return
}

func (s *ElementsStorage[T]) GetNamed(elementType T, name string) (result any, err error) {
	if elementMap, ok := s.namedElements[name]; ok && elementMap != nil {
		result, ok = elementMap[s.keyID(elementType)]
		if !ok {
			err = remyErrs.ErrElementNotRegistered{Key: elementType}
		}
		return result, err
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: elementType}
}

func (s *ElementsStorage[T]) Get(key T) (result any, err error) {
	var ok bool
	if result, ok = s.elements[s.keyID(key)]; !ok {
		err = remyErrs.ErrElementNotRegistered{Key: key}
	}
	return
}

func (s *ElementsStorage[T]) GetAll(keyTag string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		err = remyErrs.ErrConfigNotAllowReturnAll
		return
	}

	fromList := s.elements
	if keyTag != "" {
		fromList = s.namedElements[keyTag]
	}

	resultList = make([]any, 0, len(fromList))
	for _, value := range fromList {
		resultList = append(resultList, value)
	}
	return
}
