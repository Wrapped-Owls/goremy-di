package remychef

import "github.com/wrapped-owls/goremy-di/remy"

// The helpers below mirror remy's own With* family, adding the tracking that lets
// HealthCheck and Shutdown reach the instances a bind builds.

// WithLazySingleton wraps remy.RegisterLazySingleton, tracking the instance it
// builds on first retrieval.
func WithLazySingleton[T any](binder remy.Binder[T], optTag ...remy.Tag) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		remy.RegisterLazySingleton(inj, tracking(app, binder), optTag...)
	}
}

// WithSingleton wraps remy.RegisterSingleton, tracking the instance it builds
// during registration.
func WithSingleton[T any](binder remy.Binder[T], optTag ...remy.Tag) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		remy.RegisterSingleton(inj, tracking(app, binder), optTag...)
	}
}

// WithFactory wraps remy.RegisterFactory. Factory instances stay untracked: a new
// one is built on every retrieval, so keeping them would grow unbounded.
func WithFactory[T any](binder remy.Binder[T], optTag ...remy.Tag) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		remy.RegisterFactory(inj, binder, optTag...)
	}
}

// WithInstance wraps remy.RegisterInstance, tracking the already-built value.
func WithInstance[T any](value T, optTag ...remy.Tag) ModuleRegister {
	return func(app *App, inj remy.Injector) {
		app.track(value)
		remy.RegisterInstance(inj, value, optTag...)
	}
}
