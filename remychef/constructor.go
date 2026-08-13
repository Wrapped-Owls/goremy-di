package remychef

import "github.com/wrapped-owls/goremy-di/remy"

// a factory builds a new instance per retrieval, so tracking one grows unbounded
func recorderFor[T any](
	app *App,
	bindFunc func(remy.Binder[T]) remy.Bind[T],
) func(T, error) (T, error) {
	if bindFunc(nil).Type() == remy.Factory[T](nil).Type() {
		return keepResult[T]
	}

	return trackResult[T](app)
}

// WithConstructor1 registers constructor under bindFunc, resolving A from the
// injector and tracking what it builds unless bindFunc makes a factory.
func WithConstructor1[T, A any](
	bindFunc func(remy.Binder[T]) remy.Bind[T],
	constructor func(A) (T, error),
	optTag ...remy.Tag,
) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		record := recorderFor(app, bindFunc)
		remy.RegisterConstructorArgs1Err(
			inj, bindFunc,
			func(first A) (T, error) {
				return record(constructor(first))
			},
			optTag...,
		)
	}
}

// WithConstructor2 registers constructor under bindFunc, resolving A and B from
// the injector and tracking what it builds unless bindFunc makes a factory.
func WithConstructor2[T, A, B any](
	bindFunc func(remy.Binder[T]) remy.Bind[T],
	constructor func(A, B) (T, error),
	optTag ...remy.Tag,
) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		record := recorderFor(app, bindFunc)
		remy.RegisterConstructorArgs2Err(
			inj, bindFunc,
			func(first A, second B) (T, error) {
				return record(constructor(first, second))
			},
			optTag...,
		)
	}
}

// WithConstructor3 registers constructor under bindFunc, resolving A, B and C
// from the injector and tracking what it builds unless bindFunc makes a factory.
func WithConstructor3[T, A, B, C any](
	bindFunc func(remy.Binder[T]) remy.Bind[T],
	constructor func(A, B, C) (T, error),
	optTag ...remy.Tag,
) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		record := recorderFor(app, bindFunc)
		remy.RegisterConstructorArgs3Err(
			inj, bindFunc,
			func(first A, second B, third C) (T, error) {
				return record(constructor(first, second, third))
			},
			optTag...,
		)
	}
}

// WithConstructor4 registers constructor under bindFunc, resolving A, B, C and D
// from the injector and tracking what it builds unless bindFunc makes a factory.
func WithConstructor4[T, A, B, C, D any](
	bindFunc func(remy.Binder[T]) remy.Bind[T],
	constructor func(A, B, C, D) (T, error),
	optTag ...remy.Tag,
) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		record := recorderFor(app, bindFunc)
		remy.RegisterConstructorArgs4Err(
			inj, bindFunc,
			func(first A, second B, third C, fourth D) (T, error) {
				return record(constructor(first, second, third, fourth))
			},
			optTag...,
		)
	}
}
