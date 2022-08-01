package remy

import (
	"errors"
	"fmt"
	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/internal/utils"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector
	InstancePairAny     = types.InstancePair[any]

	// Bind is directly copy from types.Bind
	Bind[T any] interface {
		types.Bind[T]
	}

	// Config defines needed configuration to instantiate a new injector
	Config struct {
		// ParentInjector defines an Injector that will be used as a parent one, which will make possible to access it's
		// registered binds.
		ParentInjector Injector

		// CanOverride defines if a bind can be overridden if it is registered twice.
		CanOverride bool

		// GenerifyInterfaces defines the method to check for interface binds.
		// If this parameter is true, then an interface that is defined in two different packages,
		// but has the same signature methods, will generate the same key. If is false, all interfaces will generate
		// a different key.
		GenerifyInterfaces bool

		// UseReflectionType defines the injector to use reflection when saving and retrieving types.
		// This parameter is useful when you want to use types with different modules but the same name and package names.
		//
		// Optional, default is false.
		UseReflectionType bool
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

	reflectOpts := types.ReflectionOptions{
		GenerifyInterface: cfg.GenerifyInterfaces,
		UseReflectionType: cfg.UseReflectionType,
	}
	if cfg.ParentInjector != nil {
		return injector.New(cfg.CanOverride, reflectOpts, cfg.ParentInjector)
	}
	return injector.New(cfg.CanOverride, reflectOpts)
}

// Register must be called first, because the library doesn't support registering dependencies while get at same time.
// This is not supported in multithreading applications because it does not have race protection
func Register[T any](i Injector, bind Bind[T], keys ...string) {
	if i == nil {
		i = globalInjector()
	}
	if err := injector.Register[T](i, bind, keys...); err != nil {
		panic(err)
	}
}

// Override works like the Register function, allowing to register a bind that was already registered.
// It also must be called during binds setup, because
// the library doesn't support registering dependencies while get at same time.
//
// This is not supported in multithreading applications because it does not have race protection
func Override[T any](i Injector, bind Bind[T], keys ...string) {
	if i == nil {
		i = globalInjector()
	}
	if err := injector.Register[T](i, bind, keys...); err != nil && !errors.Is(err, utils.ErrAlreadyBound) {
		panic(err)
	}
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

// Get directly access a retriever and returns the type that was bound in it.
//
// Receives: DependencyRetriever (required); key (optional)
func Get[T any](i DependencyRetriever, keys ...string) T {
	if i == nil {
		i = globalInjector()
	}
	return injector.TryGet[T](i, keys...)
}

// DoGet directly access a retriever and returns the type that was bound in it.
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); key (optional)
func DoGet[T any](i DependencyRetriever, keys ...string) (result T, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if i == nil {
		i = globalInjector()
	}
	return injector.Get[T](i, keys...)
}

// GetGen creates a sub-injector and access the retriever to generate and return a Factory bind
//
// Receives: Injector (required); []InstancePairAny (required); key (optional)
func GetGen[T any](ij Injector, elements []InstancePairAny, keys ...string) T {
	if ij == nil {
		ij = globalInjector()
	}
	return injector.TryGetGen[T](ij, elements, keys...)
}

// DoGetGen creates a sub-injector and access the retriever to generate and return a Factory bind
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: Injector (required); []InstancePairAny (required); key (optional)
func DoGetGen[T any](ij Injector, elements []InstancePairAny, keys ...string) (result T, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if ij == nil {
		ij = globalInjector()
	}
	return injector.GetGen[T](ij, elements, keys...)
}

// GetGenFunc creates a sub-injector and access the retriever to generate and return a Factory bind
//
// Receives: Injector (required); func(Injector) (required); key (optional)
func GetGenFunc[T any](ij types.Injector, binder func(Injector), keys ...string) T {
	if ij == nil {
		ij = globalInjector()
	}
	return injector.TryGetGenFunc[T](ij, binder, keys...)
}

// DoGetGenFunc creates a sub-injector and access the retriever to generate and return a Factory bind
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: Injector (required); func(Injector) (required); key (optional)
func DoGetGenFunc[T any](ij types.Injector, binder func(Injector), keys ...string) (result T, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if ij == nil {
		ij = globalInjector()
	}
	return injector.GetGenFunc[T](ij, binder, keys...)
}
