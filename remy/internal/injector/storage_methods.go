package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

func SetStorage[T any](dStorage types.ValuesSetter[types.BindKey], value T, keys ...string) {
	var (
		key string
	)
	if len(keys) > 0 {
		key = keys[0]
	}

	if len(key) > 0 {
		dStorage.SetNamed(utils.GetKey[T](dStorage.ShouldGenerifyInterface()), key, value)
	} else {
		dStorage.Set(utils.GetKey[T](dStorage.ShouldGenerifyInterface()), value)
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
		value, ok = dStorage.GetNamed(utils.GetKey[T](dStorage.ShouldGenerifyInterface()), key)
	} else {
		value, ok = dStorage.Get(utils.GetKey[T](dStorage.ShouldGenerifyInterface()))
	}

	if ok {
		if element, assertOk := value.(T); assertOk {
			return element
		}
	}
	return utils.Default[T]()
}
