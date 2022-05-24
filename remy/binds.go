package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// Instance generates a bind that will be registered as a single instance during bind register in the Injector.
//
// This bind type has no protection over concurrency, so it's not recommended to be used a struct that performs some
// operation that will modify its attributes.
//
// It's recommended to use it only for read value injections.
// For concurrent mutable binds is recommended to use Singleton or LazySingleton
func Instance[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Instance(binder)
}

// Factory generates a bind that will be generated everytime a dependency with its type is requested by the Injector.
//
// This bind don't hold any object or pointer references, and will use the given types.Binder function everytime to
// generate a new instance of the requested type.
//
// As this bind doesn't hold any pointer and/or objects, there is no problem to use it in multiples goroutines at once.
// Just be careful with calls of the DependencyRetriever,
// if try to get and modify values from an Instance bind, it can end in a race-condition.
func Factory[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Factory(binder)
}

// Singleton generates a bind that is thread-safe, and holds the same object instance across application lifetime.
//
// The default singleton bind will execute the types.Binder during the bind registration process.
//
// If you don't want to generate the bind immediately at its registration, you can use the LazySingleton bind.
func Singleton[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Singleton(binder)
}

// LazySingleton generates a bind that works the same as the Singleton, with the only difference being that it's a lazy
// bind. So, it only will generate the singleton instance on the first Get call.
//
// It is useful in cases that you want to instantiate heavier objects only when it's needed.
func LazySingleton[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.LazySingleton(binder)
}
