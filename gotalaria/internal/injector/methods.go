package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
	"reflect"
)

func Register[T any](injector types.Injector, bind types.Bind[T]) {
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value, key := insBind.Generates(injector)
			storage.Set[T](injector.Storage(), value, key)
			return
		}
	}

	var typeT T
	elementType := reflect.TypeOf(typeT)
	injector.Bind(elementType, bind)
}

func Get[T any](injector types.DependencyRetriever) T {
	var result T
	elementType := reflect.TypeOf(result)

	// search in dynamic injections that needed to run a given function
	if bind, ok := injector.RetrieveBind(elementType); ok {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result, _ = typedBind.Generates(injector)
			return result
		}
	}

	// retrieve values from storage
	result = storage.Get[T](injector)
	return result
}
