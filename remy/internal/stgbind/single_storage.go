package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// SingleStorage is a zero-allocation storage optimized for sub-injectors that
// hold exactly one dependency. The entry is kept in flat struct fields,
// avoiding any heap allocation beyond the struct itself.
type SingleStorage struct {
	key      uint64
	name     string
	value    any
	hasValue bool
	opts     injopts.CacheConfOption
}

func NewSingleStorage(opts injopts.CacheConfOption) *SingleStorage {
	return &SingleStorage{opts: opts}
}

func (s *SingleStorage) set(key uint64, name string, value any) (wasOverridden bool, err error) {
	if s.hasValue {
		if s.key == key && s.name == name {
			if !s.opts.Is(injopts.CacheOptAllowOverride) {
				return false, remyErrs.ErrAlreadyBound{Key: types.KeyElem[any]{}}
			}
			s.value = value
			return true, nil
		}
		// A second distinct entry cannot fit in a SingleStorage.
		return false, remyErrs.ErrAlreadyBound{Key: types.KeyElem[any]{}}
	}
	s.key = key
	s.name = name
	s.value = value
	s.hasValue = true
	return false, nil
}

func (s *SingleStorage) Set(key types.BindKey, value any) (wasOverridden bool, err error) {
	return s.set(key.ID(), "", value)
}

func (s *SingleStorage) SetNamed(
	key types.BindKey, name string, value any,
) (wasOverridden bool, err error) {
	return s.set(key.ID(), name, value)
}

func (s *SingleStorage) get(key types.BindKey, name string) (result any, err error) {
	if s.hasValue && s.key == key.ID() && s.name == name {
		return s.value, nil
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: key}
}

func (s *SingleStorage) Get(key types.BindKey) (result any, err error) {
	return s.get(key, "")
}

func (s *SingleStorage) GetNamed(key types.BindKey, name string) (result any, err error) {
	return s.get(key, name)
}

func (s *SingleStorage) GetAll(keyTag string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		return nil, remyErrs.ErrConfigNotAllowReturnAll
	}
	if s.hasValue && s.name == keyTag {
		return []any{s.value}, nil
	}
	return nil, nil
}
