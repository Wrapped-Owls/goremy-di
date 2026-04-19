package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// sliceEntry is a single record held by SliceStorage.
type sliceEntry struct {
	key   uint64
	tag   string // empty string means "unnamed"
	value any
}

// SliceStorage is a flat-slice Storage optimized for a small, constant
// number of entries. It is intended for ephemeral sub-injectors (e.g., those
// created by GetWithPairs / GetWith) where the total element count is known
// up-front and is typically 1–8.
type SliceStorage struct {
	entries []sliceEntry
	opts    injopts.CacheConfOption
}

// NewSliceStorage creates a SliceStorage with the given options and an
// underlying slice pre-allocated to capacity. Pass len(elements) as
// capacity to guarantee zero re-allocations.
func NewSliceStorage(opts injopts.CacheConfOption, capacity uint) *SliceStorage {
	return &SliceStorage{
		opts:    opts,
		entries: make([]sliceEntry, 0, capacity),
	}
}

// set is the shared write path used by Set and SetNamed.
func (s *SliceStorage) set(key uint64, tag string, value any) (wasOverridden bool, err error) {
	for i := range s.entries {
		e := &s.entries[i]
		if e.key == key && e.tag == tag {
			if !s.opts.Is(injopts.CacheOptAllowOverride) {
				return false, remyErrs.ErrAlreadyBound{Key: types.KeyElem[any]{}}
			}
			e.value = value
			return true, nil
		}
	}
	s.entries = append(s.entries, sliceEntry{key: key, tag: tag, value: value})
	return false, nil
}

func (s *SliceStorage) Set(key types.BindKey, value any) (wasOverridden bool, err error) {
	return s.set(key.ID(), "", value)
}

func (s *SliceStorage) SetNamed(
	key types.BindKey, tag string, value any,
) (wasOverridden bool, err error) {
	return s.set(key.ID(), tag, value)
}

func (s *SliceStorage) get(key types.BindKey, tag string) (result any, err error) {
	id := key.ID()
	for i := range s.entries {
		e := &s.entries[i]
		if e.key == id && e.tag == tag {
			return e.value, nil
		}
	}
	return nil, remyErrs.ErrElementNotRegistered{Key: key}
}

func (s *SliceStorage) Get(key types.BindKey) (result any, err error) {
	return s.get(key, "")
}

func (s *SliceStorage) GetNamed(key types.BindKey, tag string) (result any, err error) {
	return s.get(key, tag)
}

func (s *SliceStorage) GetAll(keyTag string) (resultList []any, err error) {
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
