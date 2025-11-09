package remy

import (
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector
	InstancePairAny     = types.InstancePair[any]

	// BindKey is the internal type used to generate all type keys, and used to retrieve all types from the injector.
	// Is not supposed to use directly without the remy library, as this remove the main use of the remy-generics methods
	BindKey = types.BindKey

	// ReflectionOptions All options internally used to know how and when to use the `reflect` package
	ReflectionOptions = types.ReflectionOptions

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

		// DuckTypeElements informs to the injector that it can try to discover if the requested type
		// is on one of the already registered one.
		//
		// CAUTION: It costly a lot, since it will try to discover all registered elements
		DuckTypeElements bool

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
		GenerifyInterfaces: false,
	}
	if len(configs) > 0 {
		cfg = configs[0]
	}

	reflectOpts := types.ReflectionOptions{
		GenerifyInterface: cfg.GenerifyInterfaces,
		UseReflectionType: cfg.UseReflectionType,
	}
	cacheOpts := cacheOptsFromConfig(cfg)

	if cfg.ParentInjector != nil {
		return injector.New(cacheOpts, reflectOpts, cfg.ParentInjector)
	}
	return injector.New(cacheOpts, reflectOpts)
}

// Register must be called first, because the library doesn't support registering dependencies while get at same time.
// This is not supported in multithreading applications because it does not have race protection
func Register[T any](i Injector, bind Bind[T], keys ...string) {
	if err := injector.Register[T](mustInjector(i), bind, keys...); err != nil {
		panic(err)
	}
}

// Override works like the Register function, allowing to register a bind that was already registered.
// It also must be called during binds setup, because
// the library doesn't support registering dependencies while get at same time.
//
// This is not supported in multithreading applications because it does not have race protection
func Override[T any](i Injector, bind Bind[T], keys ...string) {
	err := injector.RegisterWithOverride[T](mustInjector(i), bind, keys...)
	if err != nil {
		panic(err)
	}
}

// RegisterInstance directly generates an instance bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterInstance[T any](i Injector, value T, keys ...string) {
	Register(mustInjector(i), Instance(value), keys...)
}

// RegisterSingleton directly generates a singleton bind without needing to write it.
//
// Receives: Injector (required); Binder (required); key (optional)
func RegisterSingleton[T any](i Injector, binder types.Binder[T], keys ...string) {
	Register(mustInjector(i), Singleton(binder), keys...)
}

// GetAll directly access a retriever and returns all instance types that was bound in it and match qualifier.
//
// Receives: DependencyRetriever (required); key (optional)
func GetAll[T any](i DependencyRetriever, keys ...string) []T {
	result, _ := injector.GetAll[T](mustRetriever(i), keys...)
	return result
}

// DoGetAll directly access a retriever and returns a list of element that match requested types that was bound in it.
// Additionally, it returns an error which indicates if the instance was found or not.
//
// Receives: DependencyRetriever (required); key (optional)
func DoGetAll[T any](i DependencyRetriever, keys ...string) (result []T, err error) {
	return injector.GetAll[T](mustRetriever(i), keys...)
}

// Get directly access a retriever and returns the type that was bound in it.
//
// Receives: DependencyRetriever (required); key (optional)
func Get[T any](i DependencyRetriever, keys ...string) T {
	return injector.TryGet[T](mustRetriever(i), keys...)
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

	return injector.Get[T](mustRetriever(i), keys...)
}

// GetGen creates a sub-injector and access the retriever to generate and return a Factory bind
//
// Receives: DependencyRetriever (required); []InstancePairAny (required); key (optional)
func GetGen[T any](i DependencyRetriever, elements []InstancePairAny, keys ...string) T {
	result, _ := injector.GetWithPairs[T](mustRetriever(i), elements, keys...)
	return result
}

// DoGetGen creates a sub-injector and access the retriever to generate and return a Factory bind
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); []InstancePairAny (required); key (optional)
func DoGetGen[T any](
	i DependencyRetriever, elements []InstancePairAny,
	keys ...string,
) (result T, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return injector.GetWithPairs[T](mustRetriever(i), elements, keys...)
}

// GetGenFunc creates a sub-injector and access the retriever to generate and return a Factory bind
//
// Receives: DependencyRetriever (required); func(Injector) (required); key (optional)
func GetGenFunc[T any](i DependencyRetriever, binder func(Injector) error, keys ...string) T {
	result, _ := injector.GetWith[T](mustRetriever(i), binder, keys...)
	return result
}

// DoGetGenFunc creates a sub-injector and access the retriever to generate and return a Factory bind
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); func(Injector) (required); key (optional)
func DoGetGenFunc[T any](
	i DependencyRetriever, binder func(Injector) error,
	keys ...string,
) (result T, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	return injector.GetWith[T](mustRetriever(i), binder, keys...)
}
