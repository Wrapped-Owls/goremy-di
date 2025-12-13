//go:build go1.24 && !nounsafe

package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

// mapBackendUnsafe implements StorageBackend[uint64, any] using a map
type mapBackendUnsafe struct {
	m genericAnyMap[uint64]
}

// newMapBackendUnsafe creates a new map backend with unsafe (uint64 keys)
func newMapBackendUnsafe() StorageBackend[uint64, any] {
	return &mapBackendUnsafe{
		m: make(genericAnyMap[uint64]),
	}
}

func (mb *mapBackendUnsafe) Set(key uint64, value any, allowOverride bool) (triedOverride bool) {
	_, triedOverride = mb.m[key]
	if triedOverride && !allowOverride {
		return true
	}
	mb.m[key] = value
	return
}

func (mb *mapBackendUnsafe) Get(key uint64) (any, error) {
	value, ok := mb.m[key]
	if !ok {
		return nil, remyErrs.ErrElementNotRegistered{Key: key}
	}
	return value, nil
}

func (mb *mapBackendUnsafe) Size() int {
	return len(mb.m)
}

func (mb *mapBackendUnsafe) GetAll() []any {
	result := make([]any, 0, len(mb.m))
	for _, v := range mb.m {
		result = append(result, v)
	}
	return result
}
