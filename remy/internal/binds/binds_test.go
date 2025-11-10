package binds

import (
	"sync"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestSingletonBind_Generates(t *testing.T) {
	var (
		expected = "gopher"
		counter  = 0
		wg       sync.WaitGroup
	)

	bind := Singleton(
		func(retriever types.DependencyRetriever) (*string, error) {
			counter += 1
			return &expected, nil
		},
	)
	// Checks if the build method is called only once
	for index := 0; index < 10; index++ {
		wg.Add(1)
		go func() {
			result, err := bind.Generates(nil)
			if err != nil {
				t.Error(err)
			}
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
