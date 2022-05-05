package binds

import "gotalaria/internal/types"

type InstanceBind[T any] struct {
	// key allows to fetch a dependency directly if it's exists.
	//
	// Optional.: default: ""
	key       string
	IsFactory bool
	binder    types.Binder[T]
}

func (b InstanceBind[T]) Generates(injector types.DependencyRetriever) (T, string) {
	return b.binder(injector), b.key
}
