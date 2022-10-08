package binds

import (
	"sync"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type SingletonBind[T any] struct {
	dependency *T
	binder     types.Binder[T]
	mutex      sync.RWMutex
	IsLazy     bool
}

func (b *SingletonBind[T]) BuildDependency(injector types.DependencyRetriever) error {
	dep, err := b.binder(injector)
	b.dependency = &dep
	return err
}

func (b *SingletonBind[T]) Generates(injector types.DependencyRetriever) (result T, err error) {
	if !b.ShouldGenerate() {
		result = *b.dependency
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Checks again if no other goroutine has initialized the dependency
	if b.dependency != nil {
		result = *b.dependency
		return
	}

	err = b.BuildDependency(injector)
	if b.dependency != nil {
		result = *b.dependency
	}

	return
}

func (b *SingletonBind[T]) Type() types.BindType {
	if b.IsLazy {
		return types.BindLazySingleton
	}
	return types.BindSingleton
}
