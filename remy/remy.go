package remy

import (
	"context"

	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector
	InstancePairAny     = types.InstancePair[any]
	BindOptions         = types.BindOptions

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

func NewBindKey[T any](_ ...T) BindKey {
	return utils.NewKeyElem[T]()
}

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
func Register[T any](i Injector, bind Bind[T], optTag ...string) {
	if err := injector.Register[T](mustInjector(i), bind, optTag...); err != nil {
		panic(err)
	}
}

// Override works like the Register function, allowing to register a bind that was already registered.
// It also must be called during binds setup, because
// the library doesn't support registering dependencies while get at same time.
//
// This is not supported in multithreading applications because it does not have race protection
func Override[T any](i Injector, bind Bind[T], optTag ...string) {
	if err := injector.RegisterWithOverride[T](mustInjector(i), bind, optTag...); err != nil {
		panic(err)
	}
}

// RegisterInstance directly generates an instance bind without needing to write it.
//
// Receives: Injector (required); value (required); tag (optional)
func RegisterInstance[T any](i Injector, value T, optTag ...string) {
	Register(mustInjector(i), Instance(value), optTag...)
}

// RegisterFactory directly generates a factory bind without needing to write it.
//
// Receives: Injector (required); Binder (required); tag (optional)
func RegisterFactory[T any](i Injector, binder types.Binder[T], optTag ...string) {
	Register(mustInjector(i), Factory(binder), optTag...)
}

// RegisterSingleton directly generates a singleton bind without needing to write it.
//
// Receives: Injector (required); Binder (required); tag (optional)
func RegisterSingleton[T any](i Injector, binder types.Binder[T], optTag ...string) {
	Register(mustInjector(i), Singleton(binder), optTag...)
}

// RegisterLazySingleton directly generates a lazy-singleton bind without needing to write it.
//
// Receives: Injector (required); Binder (required); tag (optional)
func RegisterLazySingleton[T any](i Injector, binder types.Binder[T], optTag ...string) {
	Register(mustInjector(i), LazySingleton(binder), optTag...)
}

// GetAll directly access a retriever and returns a list of element that match requested types that was bound in it.
// Additionally, it returns an error which indicates if the instance was found or not.
//
// Receives: DependencyRetriever (required); tag (optional)
func GetAll[T any](i DependencyRetriever, optTag ...string) (result []T, err error) {
	return injector.GetAll[T](mustRetriever(i), optTag...)
}

// MustGetAll directly access a retriever and returns all instance types that was bound in it and match qualifier.
// Panics if an error occurs.
//
// Receives: DependencyRetriever (required); tag (optional)
func MustGetAll[T any](i DependencyRetriever, optTag ...string) []T {
	result, err := GetAll[T](i, optTag...)
	if err != nil {
		panic(err)
	}
	return result
}

// MaybeGetAll directly access a retriever and returns all instance types that was bound in it and match qualifier.
// Returns an empty slice if an error occurs.
//
// Receives: DependencyRetriever (required); tag (optional)
func MaybeGetAll[T any](i DependencyRetriever, optTag ...string) []T {
	result, _ := GetAll[T](i, optTag...)
	return result
}

// Get directly access a retriever and returns the type that was bound in it.
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); tag (optional)
func Get[T any](i DependencyRetriever, optTag ...string) (result T, err error) {
	defer recoverInjectorPanic(&err)
	result, err = injector.Get[T](mustRetriever(i), optTag...)
	return result, err
}

// MustGet directly access a retriever and returns the type that was bound in it.
// Panics if an error occurs.
//
// Receives: DependencyRetriever (required); tag (optional)
func MustGet[T any](i DependencyRetriever, optTag ...string) T {
	result, err := Get[T](i, optTag...)
	if err != nil {
		panic(err)
	}
	return result
}

// MaybeGet directly access a retriever and returns the type that was bound in it.
// Returns the zero value of the type if an error occurs.
//
// Receives: DependencyRetriever (required); tag (optional)
func MaybeGet[T any](i DependencyRetriever, optTag ...string) T {
	result, _ := Get[T](i, optTag...)
	return result
}

// GetWithPairs creates a sub-injector and access the retriever to generate and return a Factory bind.
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); []InstancePairAny (required); tag (optional)
func GetWithPairs[T any](
	i DependencyRetriever, elements []InstancePairAny, optTag ...string,
) (result T, err error) {
	defer recoverInjectorPanic(&err)
	result, err = injector.GetWithPairs[T](mustRetriever(i), elements, optTag...)
	return result, err
}

// MustGetWithPairs creates a sub-injector and access the retriever to generate and return a Factory bind.
// Panics if an error occurs.
//
// Receives: DependencyRetriever (required); []InstancePairAny (required); tag (optional)
func MustGetWithPairs[T any](
	i DependencyRetriever, elements []InstancePairAny, optTag ...string,
) T {
	result, err := GetWithPairs[T](i, elements, optTag...)
	if err != nil {
		panic(err)
	}
	return result
}

// MaybeGetWithPairs creates a sub-injector and access the retriever to generate and return a Factory bind.
// Returns the zero value of the type if an error occurs.
//
// Receives: DependencyRetriever (required); []InstancePairAny (required); tag (optional)
func MaybeGetWithPairs[T any](
	i DependencyRetriever, elements []InstancePairAny, optTag ...string,
) T {
	result, _ := GetWithPairs[T](i, elements, optTag...)
	return result
}

// GetWith creates a sub-injector and access the retriever to generate and return a Factory bind.
// Additionally, it returns an error which indicates if the bind was found or not.
//
// Receives: DependencyRetriever (required); func(Injector) (required); tag (optional)
func GetWith[T any](
	i DependencyRetriever, binder func(Injector) error, optTag ...string,
) (result T, err error) {
	defer recoverInjectorPanic(&err)
	result, err = injector.GetWith[T](mustRetriever(i), binder, optTag...)
	return result, err
}

// MustGetWith creates a sub-injector and access the retriever to generate and return a Factory bind.
// Panics if an error occurs.
//
// Receives: DependencyRetriever (required); func(Injector) (required); tag (optional)
func MustGetWith[T any](i DependencyRetriever, binder func(Injector) error, optTag ...string) T {
	result, err := GetWith[T](i, binder, optTag...)
	if err != nil {
		panic(err)
	}
	return result
}

// MaybeGetWith creates a sub-injector and access the retriever to generate and return a Factory bind.
// Returns the zero value of the type if an error occurs.
//
// Receives: DependencyRetriever (required); func(Injector) (required); tag (optional)
func MaybeGetWith[T any](i DependencyRetriever, binder func(Injector) error, optTag ...string) T {
	result, _ := GetWith[T](i, binder, optTag...)
	return result
}

func GetWithContext[T any](
	i DependencyRetriever, ctx context.Context, optTag ...string,
) (result T, err error) {
	defer recoverInjectorPanic(&err)
	result, err = injector.GetWithPairs[T](
		i, []InstancePairAny{{Key: NewBindKey[context.Context](), Value: ctx}}, optTag...,
	)
	return result, err
}
