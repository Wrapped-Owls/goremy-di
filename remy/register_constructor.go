package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	// ConstructorEmpty defines a constructor function with no arguments that returns a value of type T and an error.
	ConstructorEmpty[T any] func() (T, error)

	// ConstructorArg1 defines a constructor function with one argument of type A that returns a value of type T and an error.
	ConstructorArg1[T, A any] func(A) (T, error)

	// ConstructorArg2 defines a constructor function with two arguments of types A and B that returns a value of type T and an error.
	ConstructorArg2[T, A, B any] func(A, B) (T, error)

	// ConstructorArg3 defines a constructor function with three arguments of types A, B, and C that returns a value of type T and an error.
	ConstructorArg3[T, A, B, C any] func(A, B, C) (T, error)

	// ConstructorArg4 defines a constructor function with four arguments of types A, B, C, and D that returns a value of type T and an error.
	ConstructorArg4[T, A, B, C, D any] func(A, B, C, D) (T, error)
)

// Binder calls the constructor function for ConstructorEmpty and returns the constructed value and any error encountered.
func (cons ConstructorEmpty[T]) Binder(types.DependencyRetriever) (T, error) {
	return cons()
}

// Binder retrieves the dependency of type A, then calls the constructor function for ConstructorArg1 and returns the constructed value and any error encountered.
func (cons ConstructorArg1[T, A]) Binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first A
	)
	if first, err = DoGet[A](retriever); err != nil {
		return
	}
	return cons(first)
}

// Binder retrieves the dependencies of types A and B, then calls the constructor function for ConstructorArg2 and returns the constructed value and any error encountered.
func (cons ConstructorArg2[T, A, B]) Binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first  A
		second B
	)
	if first, err = DoGet[A](retriever); err != nil {
		return
	}
	if second, err = DoGet[B](retriever); err != nil {
		return
	}

	return cons(first, second)
}

// Binder resolves the dependencies of the types A, B, C using the provided retriever and then calls the constructor function with these dependencies.
func (cons ConstructorArg3[T, A, B, C]) Binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first  A
		second B
		third  C
	)
	if first, err = DoGet[A](retriever); err != nil {
		return
	}
	if second, err = DoGet[B](retriever); err != nil {
		return
	}
	if third, err = DoGet[C](retriever); err != nil {
		return
	}

	return cons(first, second, third)
}

// Binder resolves the dependencies of the types A, B, C, D using the provided retriever and then calls the constructor function with these dependencies.
func (cons ConstructorArg4[T, A, B, C, D]) Binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first  A
		second B
		third  C
		fourth D
	)
	if first, err = DoGet[A](retriever); err != nil {
		return
	}
	if second, err = DoGet[B](retriever); err != nil {
		return
	}
	if third, err = DoGet[C](retriever); err != nil {
		return
	}
	if fourth, err = DoGet[D](retriever); err != nil {
		return
	}

	return cons(first, second, third, fourth)
}

// RegisterConstructorErr registers a constructor function that returns a value of type T and an error.
func RegisterConstructorErr[T any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() (T, error), keys ...string,
) {
	var generator = ConstructorEmpty[T](constructor)
	Register(mustInjector(i), bindFunc(generator.Binder), keys...)
}

// RegisterConstructor registers a constructor function that returns a value of type T without an error.
func RegisterConstructor[T any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() T, keys ...string,
) {
	var generator ConstructorEmpty[T] = func() (T, error) {
		return constructor(), nil
	}
	RegisterConstructorErr(i, bindFunc, generator, keys...)
}

// RegisterConstructorArgs1Err registers a constructor function with one argument that returns a value of type T and an error.
func RegisterConstructorArgs1Err[T, A any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A) (T, error), keys ...string,
) {
	generator := ConstructorArg1[T, A](constructor)
	Register(mustInjector(i), bindFunc(generator.Binder), keys...)
}

// RegisterConstructorArgs1 registers a constructor function with one argument that returns a value of type T without an error.
func RegisterConstructorArgs1[T, A any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A) T, keys ...string,
) {
	generator := func(arg A) (T, error) {
		return constructor(arg), nil
	}
	RegisterConstructorArgs1Err(i, bindFunc, generator, keys...)
}

// RegisterConstructorArgs2Err registers a constructor function with two arguments that returns a value of type T and an error.
func RegisterConstructorArgs2Err[T, A, B any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B) (T, error), keys ...string,
) {
	generator := ConstructorArg2[T, A, B](constructor)
	Register(mustInjector(i), bindFunc(generator.Binder), keys...)
}

// RegisterConstructorArgs2 registers a constructor function with two arguments that returns a value of type T without an error.
func RegisterConstructorArgs2[T, A, B any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B) T, keys ...string,
) {
	generator := func(arg1 A, arg2 B) (T, error) {
		return constructor(arg1, arg2), nil
	}
	RegisterConstructorArgs2Err(i, bindFunc, generator, keys...)
}

func RegisterConstructorArgs3Err[T, A, B, C any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C) (T, error), keys ...string,
) {
	generator := ConstructorArg3[T, A, B, C](constructor)
	Register(mustInjector(i), bindFunc(generator.Binder), keys...)
}

func RegisterConstructorArgs3[T, A, B, C any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C) T, keys ...string,
) {
	generator := func(arg1 A, arg2 B, arg3 C) (T, error) {
		return constructor(arg1, arg2, arg3), nil
	}
	RegisterConstructorArgs3Err(i, bindFunc, generator, keys...)
}

func RegisterConstructorArgs4Err[T, A, B, C, D any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C, D) (T, error), keys ...string,
) {
	generator := ConstructorArg4[T, A, B, C, D](constructor)
	Register(mustInjector(i), bindFunc(generator.Binder), keys...)
}

func RegisterConstructorArgs4[T, A, B, C, D any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(A, B, C, D) T, keys ...string,
) {
	generator := func(arg1 A, arg2 B, arg3 C, arg4 D) (T, error) {
		return constructor(arg1, arg2, arg3, arg4), nil
	}
	RegisterConstructorArgs4Err(i, bindFunc, generator, keys...)
}
