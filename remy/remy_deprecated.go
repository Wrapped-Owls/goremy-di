package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
)

// Deprecated: Use GetAll instead.
func DoGetAll[T any](i DependencyRetriever, optTag ...string) (result []T, err error) {
	return injector.GetAll[T](mustRetriever(i), optTag...)
}

// Deprecated: Use Get instead.
func DoGet[T any](i DependencyRetriever, optTag ...string) (result T, err error) {
	return injector.Get[T](mustRetriever(i), optTag...)
}

// Deprecated: Use MustGetGen instead.
func GetGen[T any](i DependencyRetriever, elements []InstancePairAny, optTag ...string) T {
	return MustGetWithPairs[T](i, elements, optTag...)
}

// Deprecated: Use GetGen instead.
func DoGetGen[T any](
	i DependencyRetriever, elements []InstancePairAny,
	optTag ...string,
) (result T, err error) {
	return injector.GetWithPairs[T](mustRetriever(i), elements, optTag...)
}

// Deprecated: Use MustGetGenFunc instead.
func GetGenFunc[T any](i DependencyRetriever, binder func(Injector) error, optTag ...string) T {
	return MustGetWith[T](i, binder, optTag...)
}

// Deprecated: Use GetGenFunc instead.
func DoGetGenFunc[T any](
	i DependencyRetriever, binder func(Injector) error, optTag ...string,
) (result T, err error) {
	return injector.GetWith[T](mustRetriever(i), binder, optTag...)
}
