package binds

import "github.com/wrapped-owls/fitpiece/gotalaria/internal/types"

type InstanceBind[T any] struct {
	// key allows to fetch a dependency directly if it's exists.
	//
	// Optional.: default: ""
	IsFactory bool
	binder    types.Binder[T]
}

func (b InstanceBind[T]) Generates(injector types.DependencyRetriever) T {
	return b.binder(injector)
}
