package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

var (
	__globalInjector = NewInjector()
)

// mustInjector receives an injector instance, and then returns it if exists or the global injector if it doesn't exist.
func mustInjector(ij Injector) Injector {
	if ij != nil {
		return ij
	}
	return __globalInjector
}

func mustRetriever(retriever DependencyRetriever) DependencyRetriever {
	if retriever != nil {
		return retriever
	}
	return __globalInjector
}

// SetGlobalInjector receives a custom injector and saves it to be used as a global injector
func SetGlobalInjector(i Injector) {
	if i == nil {
		panic(utils.ErrOverrideInRuntime)
	}
	__globalInjector = i
}
