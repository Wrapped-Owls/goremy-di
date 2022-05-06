package binds

import (
	"gotalaria/internal/types"
	"sync"
	"testing"
)

func TestSingletonBind_Generates(t *testing.T) {
	var (
		expected = "gopher"
		wg       sync.WaitGroup
	)

	bind := Singleton[*string](
		func(retriever types.DependencyRetriever) *string {
			return &expected
		},
	)
	// Checks if the build method is called only once
	for index := 0; index < 1000; index++ {
		wg.Add(1)
		go func() {
			result := bind.Generates(nil)
			if result != &expected {
				t.Fail()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
