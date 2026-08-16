package binds

import (
	"sync"
	"sync/atomic"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type SingletonBind[T any] struct {
	elemType[T]
	dependency atomic.Pointer[T]
	binder     types.Binder[T]
	mutex      sync.Mutex
	IsLazy     bool
}

func (b *SingletonBind[T]) BuildDependency(injector types.DependencyRetriever) error {
	dep, err := b.binder(injector)
	if err != nil {
		return err
	}

	b.dependency.Store(&dep)
	return nil
}

func (b *SingletonBind[T]) Generates(injector types.DependencyRetriever) (result T, err error) {
	if built := b.dependency.Load(); built != nil {
		return *built, nil
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Checks again if no other goroutine has initialized the dependency
	if built := b.dependency.Load(); built != nil {
		return *built, nil
	}

	if err = b.BuildDependency(injector); err != nil {
		return result, err
	}

	return *b.dependency.Load(), nil
}

func (b *SingletonBind[T]) Type() types.BindType {
	if b.IsLazy {
		return types.BindLazySingleton
	}
	return types.BindSingleton
}

func (b *SingletonBind[T]) GenAsAny(injector types.DependencyRetriever) (any, error) {
	return b.Generates(injector)
}
