package injector

import (
	"github.com/wrapped-owls/talaria-di/gotalaria/internal/types"
	"github.com/wrapped-owls/talaria-di/gotalaria/internal/utils"
)

func SetStorage[T any](dStorage types.Storage[types.BindKey], value T, keys ...string) {
	var (
		key string
	)
	if len(keys) > 0 {
		key = keys[0]
	}

	if len(key) > 0 {
		dStorage.SetNamed(utils.GetKey[T](), key, value)
	} else {
		dStorage.Set(utils.GetKey[T](), value)
	}
}

func GetStorage[T any](dStorage types.ValuesGetter[types.BindKey], keys ...string) T {
	var (
		key   string
		value any
		ok    bool
	)

	if len(keys) > 0 {
		key = keys[0]
	}

	// search in named elements
	if len(key) > 0 {
		value, ok = dStorage.GetNamed(utils.GetKey[T](), key)
	} else {
		value, ok = dStorage.Get(utils.GetKey[T]())
	}

	if ok {
		if element, assertOk := value.(T); assertOk {
			return element
		}
	}
	return utils.Default[T]()
}
