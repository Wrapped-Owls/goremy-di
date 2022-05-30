package remy

import "sync"

var (
	__globalInjector Injector
	globalMutex      sync.RWMutex
)

// globalInjector retrieves the global injector or generate a new one if it doesn't exist.
//
// This is made to prevent loading it unnecessarily into the memory
func globalInjector() Injector {
	globalMutex.RLock()
	if __globalInjector != nil {
		defer globalMutex.RUnlock()
		return __globalInjector
	}
	globalMutex.RUnlock()

	globalMutex.Lock()
	defer globalMutex.Unlock()

	// Checks again if no other goroutine has initialized the injector
	if __globalInjector == nil {
		__globalInjector = NewInjector()
	}
	return __globalInjector
}

// SetGlobalInjector receives a custom injector and saves it to be used as a global injector
func SetGlobalInjector(i Injector) {
	__globalInjector = i
}
