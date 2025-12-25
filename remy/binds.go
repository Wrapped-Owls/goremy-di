package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	// BindKey is the internal type used to generate all type keys, and used to retrieve all types from the injector.
	// Is not supposed to use directly without the remy library, as this remove the main use of the remy-generics methods
	BindKey = types.BindKey

	// Bind is directly copy from types.Bind
	Bind[T any] interface {
		types.Bind[T]
	}
)

// Instance generates a bind that will be registered as a single instance during bind register in the Injector.
//
// This bind type has no protection over concurrency, so it's not recommended to be used a struct that performs some
// operation that will modify its attributes.
//
// It's recommended to use it only for read value injections.
// For concurrent mutable binds is recommended to use Singleton or LazySingleton
func Instance[T any](element T) Bind[T] {
	return binds.Instance(element)
}

// Factory generates a bind that will be generated everytime a dependency with its type is requested by the Injector.
//
// This bind don't hold any object or pointer references, and will use the given types.Binder function everytime to
// generate a new instance of the requested type.
//
// As this bind doesn't hold any pointer and/or objects, there is no problem to use it in multiples goroutines at once.
// Just be careful with calls of the DependencyRetriever,
// if try to get and modify values from an Instance bind, it can end in a race-condition.
func Factory[T any](binder types.Binder[T]) Bind[T] {
	return binds.Factory(binder)
}

// Singleton generates a bind during the registration process.
// At the end of the registration, it holds the same object instance across application lifetime.
//
// If you don't want to generate the bind immediately at its registration, you can use the LazySingleton bind.
func Singleton[T any](binder types.Binder[T]) Bind[T] {
	return binds.Singleton(binder)
}

// LazySingleton generates a bind that works the same as the Singleton, with the only difference being that it's a lazy
// bind. So, it only will generate the singleton instance on the first Get call.
//
// It is useful in cases that you want to instantiate heavier objects only when it's needed.
func LazySingleton[T any](binder types.Binder[T]) Bind[T] {
	return binds.LazySingleton(binder)
}

// BindEntry is an interface used to pass temporary dependencies to GetWithPairs.
// It encapsulates a value along with its type key and optional tag for registration
// in a sub-injector during dependency retrieval.
type BindEntry = types.BindEntry

// NewBindEntry creates a BindEntry with the given value and no tag.
// This is a convenience function for creating temporary dependencies when using GetWithPairs.
//
// The type key is automatically generated from the value's type, so you don't need to
// manually specify it.
func NewBindEntry[T any](value T) BindEntry {
	return types.NewBindPair(value, "")
}

// NewBindEntryTagged creates a BindEntry with the given value and tag.
// This is useful when you need to register a temporary dependency with a specific tag
// for disambiguation when multiple instances of the same type are registered.
//
// The type key is automatically generated from the value's type, so you don't need to
// manually specify it.
func NewBindEntryTagged[T any](value T, tag string) BindEntry {
	return types.NewBindPair(value, tag)
}
