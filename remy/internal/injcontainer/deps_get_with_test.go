package injcontainer

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stdinj"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

func TestGetWith(t *testing.T) {
	const expected = "I love Go, yes this is true, as the answer 42"

	interfaceValue := fixtures.GoProgrammingLang{}
	testCases := [...]struct {
		name           string
		getGenCallback func(ij types.Injector) string
	}{
		{
			name: "GetWithPairs[string]",
			getGenCallback: func(ij types.Injector) string {
				result, _ := GetWithPairs[string](
					ij, "", StandardScope,
					types.NewBindPair(uint8(42), ""),
					types.NewBindPair("Go", "lang"),
					types.NewBindPair(true, ""),
					types.NewBindPair[fixtures.Language](interfaceValue, ""),
				)
				return result
			},
		},
		{
			name: "GetWith[string]",
			getGenCallback: func(i types.Injector) string {
				result, _ := GetWith[string](
					i, "", StandardScope, func(ij types.Injector) error {
						err := errors.Join(
							Register(ij, "", binds.Instance[uint8](42)),
							Register(ij, "lang", binds.Instance("Go")),
							Register(ij, "", binds.Instance(true)),
							Register[fixtures.Language](
								ij, "", binds.Instance[fixtures.Language](interfaceValue),
							),
						)
						return err
					},
				)

				return result
			},
		},
	}

	for _, tCase := range testCases {
		i := stdinj.New(stdinj.Options{Cache: injopts.CacheOptAllowOverride})
		_ = Register(
			i, "", binds.Factory(
				func(retriever types.DependencyRetriever) (result string, err error) {
					result = fmt.Sprintf(
						"I love %s, yes this is %v, as the answer %d",
						TryGet[string](
							retriever,
							"lang",
						),
						TryGet[bool](retriever, ""),
						TryGet[uint8](retriever, ""),
					)

					if _, err = Get[fixtures.Language](retriever, ""); err != nil {
						t.Error(err)
					}
					return
				},
			),
		)

		// register a bool bind to check if it will be replaced during parameter passing
		_ = Register(i, "", binds.Instance(false))

		t.Run(
			tCase.name, func(t *testing.T) {
				result := tCase.getGenCallback(i)

				if result != expected {
					t.Errorf(
						"The direct params was not injected correctly.\nExpected: `%s`\nReceived: `%s`",
						expected,
						result,
					)
					t.FailNow()
				}

				// Check if the binds don't exist after do the GetWithPairs
				var (
					uintResult, _ = Get[uint8](i, "")
					boolResult, _ = Get[bool](i, "")
					strResult, _  = Get[string](i, "lang")
				)
				if uintResult != 0 || boolResult || len(strResult) > 0 {
					t.Error("Parameter injection values override the original injector")
				}
			},
		)
	}
}

func TestGetWithPairs_withDirectBindKey(t *testing.T) {
	// Regular Show themed: Mordecai and Rigby work at the park at 3 PM
	const expected = "Mordecai and Rigby work at the park at 3 PM, during: 42 minutes, is weekend: true"

	i := stdinj.New(stdinj.Options{Cache: injopts.CacheOptAllowOverride})

	errFirstRegister := errors.Join(
		Register(
			i, "", binds.Factory(
				func(retriever types.DependencyRetriever) (result string, err error) {
					workTime := TryGet[time.Time](retriever, "")
					timeStr := workTime.Format("3 PM")
					result = fmt.Sprintf(
						"%s and %s work at the park at %s, during: %d minutes, is weekend: %v",
						TryGet[string](
							retriever,
							"employee1",
						),
						TryGet[string](retriever, "employee2"),
						timeStr,
						TryGet[uint8](retriever, ""),
						TryGet[bool](retriever, ""),
					)
					return
				},
			),
		),

		// register a bool bind to check if it will be replaced during parameter passing
		Register(i, "", binds.Instance(false)),
		Register(i, "", binds.Instance(time.Time{})),
	)
	if errFirstRegister != nil {
		t.Fatal(errFirstRegister)
	}

	// Test with direct BindKey provided - when Key is provided, InterfaceValue is not needed
	result, err := GetWithPairs[string](
		i, "", StandardScope,
		types.NewBindPair(uint8(42), ""),
		types.NewBindPair("Mordecai", "employee1"),
		types.NewBindPair("Rigby", "employee2"),
		types.NewBindPair(time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC), ""),
		types.NewBindPair(true, ""),
	)
	if err != nil {
		t.Errorf("GetWithPairs failed with error: %v", err)
		t.FailNow()
	}

	if result != expected {
		t.Errorf(
			"The direct params was not injected correctly.\nExpected: `%s`\nReceived: `%s`",
			expected,
			result,
		)
		t.FailNow()
	}

	// Check if the binds doesn't exist after do the GetWithPairs
	var (
		uintResult, _      = Get[uint8](i, "")
		boolResult, _      = Get[bool](i, "")
		employee1Result, _ = Get[string](i, "employee1")
		employee2Result, _ = Get[string](i, "employee2")
		timeResult, _      = Get[time.Time](i, "")
	)
	if uintResult != 0 || boolResult || len(employee1Result) > 0 ||
		len(employee2Result) > 0 || !timeResult.IsZero() {
		t.Error("Parameter injection values override the original injector")
	}
}

func TestGetGen_raiseCastError(t *testing.T) {
	var (
		i = stdinj.New(
			stdinj.Options{Cache: injopts.CacheOptAllowOverride},
		)
		interfaceValue fixtures.Language = fixtures.GoProgrammingLang{}
	)
	err := Register(
		i, "", binds.Factory(
			func(retriever types.DependencyRetriever) (result string, getErr error) {
				var lang fixtures.Language
				if lang, getErr = Get[fixtures.Language](retriever, ""); getErr == nil {
					result = lang.Kind() + " language: " + lang.Name()
				}
				return
			},
		),
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Run(
		"Correctly bind registration", func(t *testing.T) {
			_, err = GetWithPairs[string](
				i, "", StandardScope, types.NewBindPair[fixtures.Language](interfaceValue, ""),
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		},
	)

	t.Run(
		"Register pointer interface value", func(t *testing.T) {
			_, err = GetWithPairs[string](
				i, "", StandardScope,
				types.InstancePair[*fixtures.Language]{
					Key:   utils.NewKeyElem[fixtures.Language](),
					Value: &interfaceValue,
				},
			)
			if err == nil {
				t.Error("No error has returned after binding the value incorrectly")
				t.FailNow()
			}

			if !errors.Is(err, remyErrs.ErrTypeCastInRuntimeSentinel) {
				t.Errorf("Unknown error raised: `%v`\n", err)
			}
		},
	)
}

func TestGetWith_withParentDuckTyping(t *testing.T) {
	// Create a parent injector with CacheOptReturnAll enabled (allows GetAll)
	parent := stdinj.New(
		stdinj.Options{Cache: injopts.CacheOptReturnAll, Resolve: injopts.ResolveOptDuckTyping},
	)

	// Register an interface implementation in the parent injector
	langImpl := fixtures.GoProgrammingLang{}
	if err := Register(parent, "", binds.Instance(langImpl)); err != nil {
		t.Fatalf("Failed to register language implementation: %v", err)
	}

	// Use GetWith which creates a sub-injector with CacheOptNone (doesn't allow GetAll)
	// The sub-injector should be able to find the interface via duck typing by delegating to parent
	result, err := GetWith[fixtures.Language](
		parent, "", StandardScope, func(ij types.Injector) error {
			// Sub-injector doesn't need to register anything
			// It should find the interface from the parent via duck typing
			return nil
		},
	)
	if err != nil {
		t.Fatalf("GetWith failed to find interface via duck typing: %v", err)
	}

	// Verify the result is correct
	if result.Name() != langImpl.Name() {
		t.Errorf(
			"Language name mismatch. Expected: `%s`, Received: `%s`",
			langImpl.Name(), result.Name(),
		)
	}

	if result.Kind() != langImpl.Kind() {
		t.Errorf(
			"Language kind mismatch. Expected: `%s`, Received: `%s`",
			langImpl.Kind(), result.Kind(),
		)
	}
}

func TestGet_guessReturnsNestedMissingDependency(t *testing.T) {
	i := stdinj.New(
		stdinj.Options{Cache: injopts.CacheOptReturnAll, Resolve: injopts.ResolveOptDuckTyping},
	)
	if err := Register(
		i, "", binds.Factory(
			func(retriever types.DependencyRetriever) (fixtures.TestContextRepositoryImpl, error) {
				ctx, err := Get[context.Context](retriever, "")
				if err != nil {
					return fixtures.TestContextRepositoryImpl{}, err
				}

				requestID, _ := ctx.Value("requestID").(string)
				return fixtures.TestContextRepositoryImpl{RequestIDValue: requestID}, nil
			},
		),
	); err != nil {
		t.Fatal(err)
	}

	_, err := Get[fixtures.TestContextRepository](i, "")
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	var notRegistered remyErrs.ErrElementNotRegistered
	if !errors.As(err, &notRegistered) {
		t.Fatalf("expected ErrElementNotRegistered, got %v", err)
	}

	missingKey, ok := notRegistered.Key.(types.BindKey)
	if !ok {
		t.Fatalf("expected missing key to implement BindKey, got %T", notRegistered.Key)
	}

	expectedKey := utils.NewKeyElem[context.Context]()
	if missingKey.ID() != expectedKey.ID() {
		t.Fatalf("expected missing key to be `%T`, got `%T`", expectedKey, missingKey)
	}

	const injectedReqID = "req-123"
	ctx := context.WithValue(context.Background(), "requestID", injectedReqID)

	var result fixtures.TestContextRepository
	if result, err = GetWithPairs[fixtures.TestContextRepository](
		i, "", StandardScope, types.NewBindPair(ctx, ""),
	); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RequestID() != injectedReqID {
		t.Fatalf("unexpected request id: %s", result.RequestID())
	}
}
