//go:build nounsafe

package stgbind

import remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"

type mapBackendNounsafe[T stgKey] struct {
	m genericAnyMap[T]
}

func newBackend[T stgKey](initialCap uint16) StorageBackend[T, any] {
	return &mapBackendNounsafe[T]{
		m: make(genericAnyMap[T], initialCap),
	}
}

func (mb *mapBackendNounsafe[T]) Set(key T, value any, allowOverride bool) (triedOverride bool) {
	_, triedOverride = mb.m[key]
	if triedOverride && !allowOverride {
		return true
	}
	mb.m[key] = value
	return
}

func (mb *mapBackendNounsafe[T]) Get(key T) (any, error) {
	value, ok := mb.m[key]
	if !ok {
		return nil, remyErrs.ErrElementNotRegistered{Key: key}
	}
	return value, nil
}

func (mb *mapBackendNounsafe[T]) Size() int {
	return len(mb.m)
}

func (mb *mapBackendNounsafe[T]) GetAll() []any {
	result := make([]any, 0, len(mb.m))
	for _, v := range mb.m {
		result = append(result, v)
	}
	return result
}
