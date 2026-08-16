package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type decorable[T any] interface {
	decorate(decorator types.Decorator[T]) types.Bind[T]
}

type DecoratedBind[T any] struct {
	elemType[T]
	inner     types.Bind[T]
	decorator types.Decorator[T]
}

func Decorate[T any](bind types.Bind[T], decorator types.Decorator[T]) types.Bind[T] {
	if inner, ok := bind.(decorable[T]); ok {
		return inner.decorate(decorator)
	}

	return DecoratedBind[T]{inner: bind, decorator: decorator}
}

func decorateBinder[T any](
	binder types.Binder[T], decorator types.Decorator[T],
) types.Binder[T] {
	return func(injector types.DependencyRetriever) (T, error) {
		value, err := binder(injector)
		if err != nil {
			return value, err
		}

		return decorator(injector, value)
	}
}

func (b DecoratedBind[T]) Generates(injector types.DependencyRetriever) (result T, err error) {
	if result, err = b.inner.Generates(injector); err != nil {
		return result, err
	}

	return b.decorator(injector, result)
}

func (b DecoratedBind[T]) Type() types.BindType {
	return b.inner.Type()
}

func (b DecoratedBind[T]) GenAsAny(injector types.DependencyRetriever) (any, error) {
	return b.Generates(injector)
}
