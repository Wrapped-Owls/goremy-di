package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector

	// Bind is directly copy from types.Bind
	Bind[T any] interface {
		Generates(DependencyRetriever) T
	}

	// Config defines needed configuration to instantiate a new injector
	Config struct {
		CanOverride        bool
		GenerifyInterfaces bool
	}
)

func NewInjector(configs ...Config) Injector {
	cfg := Config{
		CanOverride:        false,
		GenerifyInterfaces: true,
	}
	if len(configs) > 0 {
		cfg = configs[0]
	}
	return injector.New(cfg.CanOverride, cfg.GenerifyInterfaces)
}

// Register must be called first, because the library doesn't support registering dependencies while get at same time.
// This is not supported in multithreading applications because it does not have race protection
func Register[T any](i Injector, bind types.Bind[T], keys ...string) {
	if i == nil {
		i = globalInjector()
	}
	injector.Register(i, bind, keys...)
}

// RegisterInstance directly generates an instance bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterInstance[T any](i Injector, value T, keys ...string) {
	Register(
		i,
		Instance(func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

// RegisterSingleton directly generates a singleton bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterSingleton[T any](i Injector, value T, keys ...string) {
	Register(
		i,
		Singleton(func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

func Get[T any](i DependencyRetriever, keys ...string) T {
	if i == nil {
		i = globalInjector()
	}
	return injector.Get[T](i, keys...)
}
