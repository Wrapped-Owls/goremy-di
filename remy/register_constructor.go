package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	// ConstructorEmpty defines a constructor function with no arguments that returns a value of type T and an error.
	ConstructorEmpty[T any] func() (T, error)

	// ConstructorArg1 defines a constructor function with one argument of type K that returns a value of type T and an error.
	ConstructorArg1[T, K any] func(K) (T, error)

	// ConstructorArg2 defines a constructor function with two arguments of types K and P that returns a value of type T and an error.
	ConstructorArg2[T, K, P any] func(K, P) (T, error)
)

// binder calls the constructor function for ConstructorEmpty and returns the constructed value and any error encountered.
func (cons ConstructorEmpty[T]) binder(types.DependencyRetriever) (T, error) {
	return cons()
}

// binder retrieves the dependency of type K, then calls the constructor function for ConstructorArg1 and returns the constructed value and any error encountered.
func (cons ConstructorArg1[T, K]) binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first K
	)
	if first, err = DoGet[K](retriever); err != nil {
		return
	}
	return cons(first)
}

// binder retrieves the dependencies of types K and P, then calls the constructor function for ConstructorArg2 and returns the constructed value and any error encountered.
func (cons ConstructorArg2[T, K, P]) binder(retriever types.DependencyRetriever) (value T, err error) {
	var (
		first  K
		second P
	)
	if first, err = DoGet[K](retriever); err != nil {
		return
	}
	if second, err = DoGet[P](retriever); err != nil {
		return
	}

	return cons(first, second)
}

// RegisterConstructorErr registers a constructor function that returns a value of type T and an error.
func RegisterConstructorErr[T any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func() (T, error), keys ...string,
) {
	var generator = ConstructorEmpty[T](constructor)
	Register(mustInjector(i), bindFunc(generator.binder), keys...)
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
func RegisterConstructorArgs1Err[T, K any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(K) (T, error), keys ...string,
) {
	generator := ConstructorArg1[T, K](constructor)
	Register(mustInjector(i), bindFunc(generator.binder), keys...)
}

// RegisterConstructorArgs1 registers a constructor function with one argument that returns a value of type T without an error.
func RegisterConstructorArgs1[T, K any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(K) T, keys ...string,
) {
	generator := func(arg K) (T, error) {
		return constructor(arg), nil
	}
	RegisterConstructorArgs1Err(i, bindFunc, generator, keys...)
}

// RegisterConstructorArgs2Err registers a constructor function with two arguments that returns a value of type T and an error.
func RegisterConstructorArgs2Err[T, K, P any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(K, P) (T, error), keys ...string,
) {
	generator := ConstructorArg2[T, K, P](constructor)
	Register(mustInjector(i), bindFunc(generator.binder), keys...)
}

// RegisterConstructorArgs2 registers a constructor function with two arguments that returns a value of type T without an error.
func RegisterConstructorArgs2[T, K, P any](
	i Injector, bindFunc func(binder types.Binder[T]) Bind[T],
	constructor func(K, P) T, keys ...string,
) {
	generator := func(arg1 K, arg2 P) (T, error) {
		return constructor(arg1, arg2), nil
	}
	RegisterConstructorArgs2Err(i, bindFunc, generator, keys...)
}
