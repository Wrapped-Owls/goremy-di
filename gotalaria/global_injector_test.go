package gotalaria

import (
	"sync"
	"testing"
)

func TestGlobal_GetWithConcurrency(t *testing.T) {
	const expected = "it's working, yaaay!"
	var wg sync.WaitGroup

	Register(nil, Instance(func(retriever DependencyRetriever) string {
		return expected
	}))
	for i := 0; i < 255; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if Get[string](nil) != expected {
				t.Error("string obtained is not the same as expected")
			}
		}()
	}
	wg.Wait()
}
