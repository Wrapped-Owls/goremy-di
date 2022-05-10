package gotalaria

import (
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/injector"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/types"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector

	// Bind is directly copy from types.Bind
	Bind[T any] interface {
		Generates(DependencyRetriever) T
	}
)

func NewInjector() Injector {
	return injector.New(false)
}

func Register[T any](i Injector, bind types.Bind[T], keys ...string) {
	injector.Register[T](i, bind, keys...)
}

// RegisterInstance directly generates an instance bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterInstance[T any](i Injector, value T, keys ...string) {
	injector.Register[T](
		i,
		Instance[T](func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

// RegisterSingleton directly generates a singleton bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterSingleton[T any](i Injector, value T, keys ...string) {
	injector.Register[T](
		i,
		Singleton[T](func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

func Get[T any](i DependencyRetriever, keys ...string) T {
	return injector.Get[T](i, keys...)
}
