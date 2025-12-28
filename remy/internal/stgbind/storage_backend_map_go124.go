//go:build go1.24 && !nounsafe

package stgbind

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

// mapBackendUnsafe implements StorageBackend[uint64, any] using a map
type mapBackendUnsafe[T stgKey] struct {
	m genericAnyMap[uint64]
}

// newBackend creates a new map backend with unsafe (uint64 keys)
func newBackend[T stgKey]() StorageBackend[T, any] {
	return &mapBackendUnsafe[T]{
		m: make(genericAnyMap[uint64]),
	}
}

func (mb *mapBackendUnsafe[T]) Set(key T, value any, allowOverride bool) (triedOverride bool) {
	keyID := key.ID()
	_, triedOverride = mb.m[keyID]
	if triedOverride && !allowOverride {
		return true
	}
	mb.m[keyID] = value
	return
}

func (mb *mapBackendUnsafe[T]) Get(key T) (any, error) {
	keyID := key.ID()
	value, ok := mb.m[keyID]
	if !ok {
		return nil, remyErrs.ErrElementNotRegistered{Key: key}
	}
	return value, nil
}

func (mb *mapBackendUnsafe[T]) Size() int {
	return len(mb.m)
}

func (mb *mapBackendUnsafe[T]) GetAll() []any {
	result := make([]any, 0, len(mb.m))
	for _, v := range mb.m {
		result = append(result, v)
	}
	return result
}
