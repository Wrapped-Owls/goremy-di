package binds

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"sync"
	"testing"
)

func TestSingletonBind_Generates(t *testing.T) {
	var (
		expected = "gopher"
		counter  = 0
		wg       sync.WaitGroup
	)

	bind := Singleton[*string](
		func(retriever types.DependencyRetriever) *string {
			counter += 1
			return &expected
		},
	)
	// Checks if the build method is called only once
	for index := 0; index < 10; index++ {
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
	if counter > 1 {
		t.Errorf("function `Bind.Generates` executed %d times", counter)
	}
}
