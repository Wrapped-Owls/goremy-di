package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
	"gotalaria/internal/utils"
)

func Register[T any](injector types.Injector, bind types.Bind[T], keys ...string) {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value := insBind.Generates(injector)
			storage.Set[T](injector, value, key)
			return
		}
	} else if sglBind, assertOk := bind.(*binds.SingletonBind[T]); assertOk {
		if !sglBind.IsLazy && sglBind.ShouldGenerate() {
			sglBind.BuildDependency(injector)
		}
	}

	elementType := utils.GetKey[T]()

	if len(key) > 0 {
		injector.BindNamed(key, elementType, bind)
	} else {
		injector.Bind(elementType, bind)
	}
}

func Get[T any](injector types.DependencyRetriever, keys ...string) T {
	var (
		key string
	)

	if len(keys) > 0 {
		key = keys[0]
	}
	elementType := utils.GetKey[T]()

	var (
		bind any
		ok   bool
	)

	// search in dynamic injections that needed to run a given function
	if len(key) > 0 {
		bind, ok = injector.RetrieveNamedBind(key, elementType)
	} else {
		bind, ok = injector.RetrieveBind(elementType)
	}

	if ok {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result := typedBind.Generates(injector)
			return result
		}
	}
	// retrieve values from instanceStorage
	result := storage.Get[T](injector, key)
	return result
}
