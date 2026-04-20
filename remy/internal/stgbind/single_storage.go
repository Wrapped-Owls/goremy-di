package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// SingleStorage is a zero-allocation storage optimized for sub-injectors that
// hold exactly one dependency. The entry is kept in flat struct fields,
// avoiding any heap allocation beyond the struct itself.
type SingleStorage[T stgKey] struct {
	baseStorage[T]
	key      bindKeyID
	name     string
	value    any
	hasValue bool
}

func NewSingleStorage[T stgKey](opts injopts.CacheConfOption) *SingleStorage[T] {
	return &SingleStorage[T]{
		baseStorage: newBaseStorage[T](opts),
	}
}

func (s *SingleStorage[T]) set(key T, name string, value any) (wasOverridden bool, err error) {
	checkKey := s.keyID(key)
	if s.hasValue {
		if s.key == checkKey && s.name == name {
			if !s.opts.Is(injopts.CacheOptAllowOverride) {
				return false, remyErrs.ErrAlreadyBound{Key: key}
			}
			s.value = value
			return true, nil
		}
		// A second distinct entry cannot fit in a SingleStorage.
		return false, remyErrs.ErrAlreadyBound{Key: key}
	}
	s.key = checkKey
	s.name = name
	s.value = value
	s.hasValue = true
	return false, nil
}

func (s *SingleStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	return s.set(key, "", value)
}

func (s *SingleStorage[T]) SetNamed(
	key T, name string, value any,
) (wasOverridden bool, err error) {
	return s.set(key, name, value)
}

func (s *SingleStorage[T]) get(key T, name string) (result any, err error) {
	if s.hasValue && s.key == s.keyID(key) && s.name == name {
		return s.value, nil
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: key}
}

func (s *SingleStorage[T]) Get(key T) (result any, err error) {
	return s.get(key, "")
}

func (s *SingleStorage[T]) GetNamed(key T, name string) (result any, err error) {
	return s.get(key, name)
}

func (s *SingleStorage[T]) GetAll(keyTag string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		return nil, remyErrs.ErrConfigNotAllowReturnAll
	}
	if s.hasValue && s.name == keyTag {
		return []any{s.value}, nil
	}
	return nil, nil
}
