package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

func SetStorage[T any](dStorage types.ValuesSetter[types.BindKey], value T, keys ...string) error {
	var (
		key string
	)
	if len(keys) > 0 {
		key = keys[0]
	}

	if len(key) > 0 {
		return dStorage.SetNamed(utils.GetKey[T](dStorage.ReflectOpts()), key, value)
	}
	return dStorage.Set(utils.GetKey[T](dStorage.ReflectOpts()), value)
}

func GetStorage[T any](dStorage types.ValuesGetter[types.BindKey], keys ...string) (T, error) {
	var (
		key   string
		value any
		err   error
	)

	if len(keys) > 0 {
		key = keys[0]
	}

	// search in named elements
	if len(key) > 0 {
		value, err = dStorage.GetNamed(utils.GetKey[T](dStorage.ReflectOpts()), key)
	} else {
		value, err = dStorage.Get(utils.GetKey[T](dStorage.ReflectOpts()))
	}

	if err == nil {
		if element, assertOk := value.(T); assertOk {
			return element, nil
		}
	}
	return utils.Default[T](), err
}
