package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

func Register[T any](ij types.Injector, bind types.Bind[T], keys ...string) error {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	elementType := utils.GetKey[T](ij.ReflectOpts())
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value := insBind.Generates(ij)
			if len(key) > 0 {
				return ij.BindNamed(elementType, key, value)
			}
			return ij.Bind(elementType, value)
		}
	} else if sglBind, assertOk := bind.(*binds.SingletonBind[T]); assertOk {
		if !sglBind.IsLazy && sglBind.ShouldGenerate() {
			sglBind.BuildDependency(ij)
		}
	}

	if len(key) > 0 {
		return ij.BindNamed(elementType, key, bind)
	}
	return ij.Bind(elementType, bind)
}

func Get[T any](retriever types.DependencyRetriever, keys ...string) (T, error) {
	var (
		key string
	)

	if len(keys) > 0 {
		key = keys[0]
	}
	elementType := utils.GetKey[T](retriever.ReflectOpts())

	var (
		bind any
		err  error
	)

	// search in dynamic injections that needed to run a given function
	if len(key) > 0 {
		bind, err = retriever.RetrieveNamedBind(elementType, key)
	} else {
		bind, err = retriever.RetrieveBind(elementType)
	}

	if err == nil {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result := typedBind.Generates(retriever)
			return result, nil
		}
		if instanceBind, assertOk := bind.(T); assertOk {
			return instanceBind, nil
		}
	}
	// retrieve values from cacheStorage
	return utils.Default[T](), err
}

func TryGet[T any](retriever types.DependencyRetriever, keys ...string) (result T) {
	result, _ = Get[T](retriever, keys...)
	return
}

func GetGen[T any](ij types.Injector, elements []types.InstancePair[any], keys ...string) (result T, err error) {
	subInjector := New(false, ij.ReflectOpts(), ij)
	for _, element := range elements {
		bindKey := utils.GetElemKey(element.Value, subInjector.ReflectOpts())
		if len(element.Key) > 0 {
			if err = subInjector.BindNamed(bindKey, element.Key, element.Value); err != nil {
				return
			}
		} else if err = subInjector.Bind(bindKey, element.Value); err != nil {
			return
		}
	}

	return Get[T](subInjector, keys...)
}

func TryGetGen[T any](ij types.Injector, elements []types.InstancePair[any], keys ...string) (result T) {
	result, _ = GetGen[T](ij, elements, keys...)
	return
}

func GetGenFunc[T any](ij types.Injector, binder func(injector types.Injector), keys ...string) (T, error) {
	subInjector := New(false, ij.ReflectOpts(), ij)
	binder(subInjector)
	return Get[T](subInjector, keys...)
}

func TryGetGenFunc[T any](ij types.Injector, binder func(injector types.Injector), keys ...string) (result T) {
	result, _ = GetGenFunc[T](ij, binder, keys...)
	return
}
