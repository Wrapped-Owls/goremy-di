package storage

import "gotalaria/internal/types"

func Set[T any](dStorage types.Storage, value T, keys ...string) {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	if len(key) > 0 {
		dStorage.SetNamed(key, value)
		return
	}
	dStorage.Set(value)
}

func Get[T any](dStorage types.Storage, keys ...string) T {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	// search in named binds
	if len(key) > 0 {
		if value := dStorage.Get(key); value != nil {
			if element, assertOk := value.(T); assertOk {
				return element
			}
		}
	}

	// search in unsorted binds
	for _, bind := range dStorage.Binds() {
		if result, ok := bind.(T); ok {
			return result
		}
	}

	var defaultResult T
	return defaultResult
}
