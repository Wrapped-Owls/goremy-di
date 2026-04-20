package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// sliceEntry is a single record held by SliceStorage.
type sliceEntry struct {
	key   bindKeyID
	tag   string // empty string means "unnamed"
	value any
}

// SliceStorage is a flat-slice Storage optimized for a small, constant
// number of entries. It is intended for ephemeral sub-injectors (e.g., those
// created by GetWithPairs / GetWith) where the total element count is known
// up-front and is typically 1–8.
type SliceStorage[T stgKey] struct {
	baseStorage[T]
	entries []sliceEntry
}

func NewSliceStorage[T stgKey](opts injopts.CacheConfOption, capacity uint) *SliceStorage[T] {
	return &SliceStorage[T]{
		baseStorage: newBaseStorage[T](opts),
		entries:     make([]sliceEntry, 0, capacity),
	}
}

// set is the shared write path used by Set and SetNamed.
func (s *SliceStorage[T]) set(key T, tag string, value any) (wasOverridden bool, err error) {
	checkKey := s.keyID(key)
	for i := range s.entries {
		e := &s.entries[i]
		if e.key == checkKey && e.tag == tag {
			if !s.opts.Is(injopts.CacheOptAllowOverride) {
				return false, remyErrs.ErrAlreadyBound{Key: key}
			}
			e.value = value
			return true, nil
		}
	}
	s.entries = append(s.entries, sliceEntry{key: checkKey, tag: tag, value: value})
	return false, nil
}

func (s *SliceStorage[T]) Set(key T, value any) (wasOverridden bool, err error) {
	return s.set(key, "", value)
}

func (s *SliceStorage[T]) SetNamed(
	key T, tag string, value any,
) (wasOverridden bool, err error) {
	return s.set(key, tag, value)
}

func (s *SliceStorage[T]) get(key T, tag string) (result any, err error) {
	id := s.keyID(key)
	for i := range s.entries {
		e := &s.entries[i]
		if e.key == id && e.tag == tag {
			return e.value, nil
		}
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: key}
}

func (s *SliceStorage[T]) Get(key T) (result any, err error) {
	return s.get(key, "")
}

func (s *SliceStorage[T]) GetNamed(key T, tag string) (result any, err error) {
	return s.get(key, tag)
}

func (s *SliceStorage[T]) GetAll(keyTag string) (resultList []any, err error) {
	if !s.opts.Is(injopts.CacheOptReturnAll) {
		return nil, remyErrs.ErrConfigNotAllowReturnAll
	}

	resultList = make([]any, 0, len(s.entries))
	for i := range s.entries {
		e := &s.entries[i]
		if e.tag == keyTag {
			resultList = append(resultList, e.value)
		}
	}
	return resultList, nil
}
