package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
	"reflect"
)

func Register[T any](injector types.Injector, bind types.Bind[T], keys ...string) {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value, _ := insBind.Generates(injector)
			storage.Set[T](injector.Storage(), value, key)
			return
		}
	}

	var typeT T
	elementType := reflect.TypeOf(typeT)

	if len(key) > 0 {
		injector.BindNamed(key, elementType, bind)
	} else {
		injector.Bind(elementType, bind)
	}
}

func Get[T any](injector types.DependencyRetriever, keys ...string) T {
	var (
		key    string
		result T
	)

	if len(keys) > 0 {
		key = keys[0]
	}
	elementType := reflect.TypeOf(result)

	var (
		bind any
		ok   bool
	)

	if len(key) > 0 {
		bind, ok = injector.RetrieveNamedBind(key, elementType)
	} else {
		bind, ok = injector.RetrieveBind(elementType)
	}

	// search in dynamic injections that needed to run a given function
	if ok {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result, _ = typedBind.Generates(injector)
			return result
		}
	}
	// retrieve values from storage
	result = storage.Get[T](injector, key)
	return result
}
