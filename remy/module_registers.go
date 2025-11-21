package remy

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

// The following helpers adapt the existing registering API (which requires an Injector)
// into ModuleRegister functions that can be used with NewModule. They DO NOT change
// the current API; they are only convenience wrappers for module composition.

// WithBind wraps Register into a ModuleRegister.
func WithBind[T any](bind Bind[T], optTag ...string) ModuleRegister {
	return func(i Injector) { Register(i, bind, optTag...) }
}

// WithInstance wraps RegisterInstance into a ModuleRegister.
func WithInstance[T any](value T, optTag ...string) ModuleRegister {
	return func(i Injector) { RegisterInstance(i, value, optTag...) }
}

// WithFactory wraps RegisterFactory into a ModuleRegister.
func WithFactory[T any](binder types.Binder[T], optTag ...string) ModuleRegister {
	return func(i Injector) { RegisterFactory(i, binder, optTag...) }
}

// WithSingleton wraps RegisterSingleton into a ModuleRegister.
func WithSingleton[T any](binder types.Binder[T], optTag ...string) ModuleRegister {
	return func(i Injector) { RegisterSingleton(i, binder, optTag...) }
}

// WithLazySingleton wraps RegisterLazySingleton into a ModuleRegister.
func WithLazySingleton[T any](binder types.Binder[T], optTag ...string) ModuleRegister {
	return func(i Injector) { RegisterLazySingleton(i, binder, optTag...) }
}

// WithConstructor and its arity variants wrap the constructor helpers for module composition.

func WithConstructor[T any](
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() (T, error), optTag ...string,
) ModuleRegister {
	return func(i Injector) { RegisterConstructorErr(i, bindFunc, constructor, optTag...) }
}

func WithConstructor1[T any, A any](
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A) (T, error), optTag ...string,
) ModuleRegister {
	return func(i Injector) { RegisterConstructorArgs1Err(i, bindFunc, constructor, optTag...) }
}

func WithConstructor2[T any, A, B any](
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B) (T, error), optTag ...string,
) ModuleRegister {
	return func(i Injector) { RegisterConstructorArgs2Err(i, bindFunc, constructor, optTag...) }
}

func WithConstructor3[T any, A, B, C any](
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C) (T, error), optTag ...string,
) ModuleRegister {
	return func(i Injector) { RegisterConstructorArgs3Err(i, bindFunc, constructor, optTag...) }
}

func WithConstructor4[T any, A, B, C, D any](
	bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C, D) (T, error), optTag ...string,
) ModuleRegister {
	return func(i Injector) { RegisterConstructorArgs4Err(i, bindFunc, constructor, optTag...) }
}
