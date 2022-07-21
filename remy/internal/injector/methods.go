package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

func Register[T any](ij types.Injector, bind types.Bind[T], keys ...string) {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value := insBind.Generates(ij)
			SetStorage(ij, value, key)
			return
		}
	} else if sglBind, assertOk := bind.(*binds.SingletonBind[T]); assertOk {
		if !sglBind.IsLazy && sglBind.ShouldGenerate() {
			sglBind.BuildDependency(ij)
		}
	}

	elementType := utils.GetKey[T](ij.ReflectOpts())

	if len(key) > 0 {
		ij.BindNamed(key, elementType, bind)
	} else {
		ij.Bind(elementType, bind)
	}
}

func Get[T any](retriever types.DependencyRetriever, keys ...string) T {
	var (
		key string
	)

	if len(keys) > 0 {
		key = keys[0]
	}
	elementType := utils.GetKey[T](retriever.ReflectOpts())

	var (
		bind any
		ok   bool
	)

	// search in dynamic injections that needed to run a given function
	if len(key) > 0 {
		bind, ok = retriever.RetrieveNamedBind(key, elementType)
	} else {
		bind, ok = retriever.RetrieveBind(elementType)
	}

	if ok {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result := typedBind.Generates(retriever)
			return result
		}
	}
	// retrieve values from instanceStorage
	result := GetStorage[T](retriever, key)
	return result
}

func GetGen[T any](ij types.Injector, elements []types.InstancePair[any], keys ...string) T {
	subInjector := New(false, ij.ReflectOpts(), ij)
	for _, element := range elements {
		bindKey := utils.GetElemKey(element.Value, subInjector.ReflectOpts())
		if len(element.Key) > 0 {
			subInjector.SetNamed(bindKey, element.Key, element.Value)
		}
		subInjector.Set(bindKey, element.Value)
	}

	return Get[T](subInjector, keys...)
}

func GetGenFunc[T any](ij types.Injector, binder func(injector types.Injector), keys ...string) T {
	subInjector := New(false, ij.ReflectOpts(), ij)
	binder(subInjector)
	return Get[T](subInjector, keys...)
}
