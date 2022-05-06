package binds

import (
	"gotalaria/internal/types"
	"sync"
)

type SingletonBind[T any] struct {
	dependency *T
	binder     types.Binder[T]
	mutex      sync.Mutex
	IsLazy     bool
}

func (b *SingletonBind[T]) BuildDependency(injector types.DependencyRetriever) {
	dep := b.binder(injector)
	b.dependency = &dep
}

func (b *SingletonBind[T]) Generates(injector types.DependencyRetriever) T {
	b.mutex.Lock()
	if b.dependency == nil {
		b.BuildDependency(injector)
	}
	b.mutex.Unlock()

	element := *b.dependency
	return element
}
