package injector

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

// Test the thread-safety of Get, GetAll, GetWith, and GetWithPairs by registering
// everything on the main goroutine and performing concurrent reads from workers.
func TestMethods_parallel_Get_variants(t *testing.T) {
	const (
		workers    = 32
		iterations = 1000
	)

	// Create injector and register all dependencies on the main goroutine
	i := New(injopts.CacheOptReturnAll)

	if registerErr := errors.Join(
		// Base registrations that should remain unchanged across all workers
		Register(i, binds.Instance(uint8(42))),
		Register(i, binds.Instance("Go"), "lang"),
		Register(i, binds.Instance(false)), // will be overridden only in sub-injectors
		Register[fixtures.Language](
			i, binds.Instance[fixtures.Language](fixtures.GoProgrammingLang{}),
		),
	); registerErr != nil {
		t.Fatal(registerErr)
	}

	// Register a factory that depends on values (uint8, bool, string[tag], interface)
	if err := Register(i, binds.Factory(func(retriever types.DependencyRetriever) (string, error) {
		// Compose a stable string:
		// "I love <lang>, yes this is <bool>, as the answer <uint8>"
		res := fmt.Sprintf(
			"I love %s, yes this is %v, as the answer %d",
			TryGet[string](retriever, "lang"),
			TryGet[bool](retriever),
			TryGet[uint8](retriever),
		)

		// Make sure interface retrieval also works
		if _, err := Get[fixtures.Language](retriever); err != nil {
			return "", err
		}
		return res, nil
	})); err != nil {
		t.Fatalf("failed to register string factory: %v", err)
	}

	const expectedStr = "I love Go, yes this is true, as the answer 42"

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for it := 0; it < iterations; it++ {
				// Simple Get
				if v, err := Get[uint8](i); err != nil || v != 42 {
					t.Errorf("failed on Get[uint8]: val=%d err=%v", v, err)
					return
				}
				if s, err := Get[string](i, "lang"); err != nil || s != "Go" {
					t.Errorf("failed on Get[string](lang): val=%s err=%v", s, err)
					return
				}

				// GetAll
				if lst, err := GetAll[uint8](i); err != nil || len(lst) != 1 || lst[0] != 42 {
					t.Errorf("failed on GetAll[uint8]: list=%v err=%v", lst, err)
					return
				}

				// GetWithPairs overriding only within the sub-injector
				valPairs, err := GetWithPairs[string](
					i,
					[]types.InstancePair[any]{
						{
							Value: true, // override bool only for this call
							Key:   utils.NewKeyElem[bool](),
						},
						{
							Value: fixtures.CountryLanguage{Language: "ptBr"},
							Key:   utils.NewKeyElem[fixtures.Language](),
						},
					},
				)
				if err != nil || valPairs != expectedStr {
					t.Errorf("failed on GetWithPairs[string]: val=%q err=%v", valPairs, err)
					return
				}

				// Ensure original bool remains false in the main injector
				if b, getErr := Get[bool](i); getErr != nil || b != false {
					t.Errorf("post-GetWithPairs Get[bool]: val=%v err=%v", b, getErr)
					return
				}

				// GetWith binder variant
				var valWith string
				valWith, err = GetWith[string](i, func(ij types.Injector) error {
					// supply only the overrides for this call
					return errors.Join(
						Register(ij, binds.Instance(true)),
						Register(
							ij, binds.Instance[fixtures.Language](
								fixtures.CountryLanguage{Language: "ptBr"},
							),
						),
					)
				})
				if err != nil || valWith != expectedStr {
					t.Errorf("failed on GetWith[string]: val=%q err=%v", valWith, err)
					return
				}
			}
		}()
	}

	wg.Wait()
}
