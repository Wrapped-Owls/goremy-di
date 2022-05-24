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
			if len(Get[string](nil)) != 0 {
				t.Error("string retrieved from global injector should be empty")
			}
		}()
	}
	wg.Wait()
}
