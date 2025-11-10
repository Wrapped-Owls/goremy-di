package remy

import (
	"sync"
	"testing"
)

func TestGlobal_GetWithConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 255; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(MaybeGet[string](nil)) != 0 {
				t.Error("string retrieved from global injector should be empty")
			}
		}()
	}
	wg.Wait()
}

func TestSetGlobalInjector(t *testing.T) {
	ij := NewInjector()
	var counter uint8 = 0
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (uint8, error) {
				counter += 1
				return counter, nil
			},
		),
	)

	value := MustGet[uint8](ij)
	if value != 1 {
		t.Errorf("value should be 1, got %d", value)
	}
	SetGlobalInjector(ij)
	value = MustGet[uint8](nil)
	if value != 2 {
		t.Errorf("value should be 2, got %d", value)
	}
}
