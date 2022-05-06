package binds

import (
	"gotalaria/internal/types"
	"sync"
)

type SingletonBind[T any] struct {
	dependency *T
	binder     types.Binder[T]
	mutex      sync.RWMutex
	IsLazy     bool
}

func (b *SingletonBind[T]) BuildDependency(injector types.DependencyRetriever) {
	dep := b.binder(injector)
	b.dependency = &dep
}

func (b *SingletonBind[T]) Generates(injector types.DependencyRetriever) T {
	b.mutex.RLock()
	if b.dependency != nil {
		defer b.mutex.RUnlock()
		return *b.dependency
	}
	b.mutex.RUnlock()

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Checks again if no other goroutine has initialized the dependency
	if b.dependency != nil {
		return *b.dependency
	}
	b.BuildDependency(injector)

	return *b.dependency
}
