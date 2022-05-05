package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/storage"
	"gotalaria/internal/types"
	"reflect"
)

func Register[T any](injector *StdInjector, bind types.Bind[T]) {
	if insBind, ok := bind.(binds.InstanceBind[T]); ok {
		if !insBind.IsFactory {
			value, key := insBind.Generates(injector)
			storage.Set[T](injector.storage, value, key)
			return
		}
	}

	var typeT T
	elementType := reflect.TypeOf(typeT)
	injector.dynamicDependencies[elementType] = bind
}

func Get[T any](injector *StdInjector) T {
	var result T
	elementType := reflect.TypeOf(result)

	if bind, ok := injector.dynamicDependencies[elementType]; ok {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			result, _ = typedBind.Generates(injector)
			return result
		}
	}
	return result
}
