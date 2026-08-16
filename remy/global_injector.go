package remy

import (
	"errors"
)

var (
	__globalInjector               = NewInjector()
	ErrOverrideGlobalInjectWithNil = errors.New("tried to override global injector with nil param")
)

// injectorOrGlobal returns the given injector, falling back to the global one when it is nil.
func injectorOrGlobal(ij Injector) Injector {
	if ij != nil {
		return ij
	}
	return __globalInjector
}

func retrieverOrGlobal(retriever DependencyRetriever) DependencyRetriever {
	if retriever != nil {
		return retriever
	}
	return __globalInjector
}

// SetGlobalInjector receives a custom injector and saves it to be used as a global injector
func SetGlobalInjector(i Injector) {
	if i == nil {
		panic(ErrOverrideGlobalInjectWithNil)
	}
	__globalInjector = i
}
