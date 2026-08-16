package injcontainer

import (
	"errors"
	"runtime"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stdinj"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

func TestGet_duckTypeInterface(t *testing.T) {
	strGenerator := func(lang fixtures.Language) string {
		return lang.Kind() + " language: " + lang.Name()
	}

	var (
		testFirstSubject  = fixtures.GoProgrammingLang{}
		testSecondSubject = fixtures.CountryLanguage{}
		testCases         = [...]struct {
			name            string
			registerSubject uint8
			expected        string
			expectedError   error
		}{
			{
				name:            "Correctly bind registration",
				registerSubject: 1,
				expected:        strGenerator(testFirstSubject),
			},
			{
				name:            "Failed to find dependency bind",
				registerSubject: 0,
				expected:        "",
				expectedError:   remyErrs.ErrElementNotRegisteredSentinel,
			},
			{
				name:            "Inject multiple elements that implements interface",
				registerSubject: 2,
				expected:        "",
				expectedError:   remyErrs.ErrFoundMoreThanOneValidDISentinel,
			},
		}
	)

	for _, tt := range testCases {
		t.Run(
			tt.name, func(t *testing.T) {
				i := stdinj.New(
					stdinj.Options{
						Cache:   injopts.CacheOptReturnAll,
						Resolve: injopts.ResolveOptDuckTyping,
					},
				)
				err := Register(
					i, "", binds.Factory(
						func(retriever types.DependencyRetriever) (result string, getErr error) {
							var lang fixtures.Language
							if lang, getErr = Get[fixtures.Language](retriever, ""); getErr == nil {
								result = strGenerator(lang)
							}
							return
						},
					),
				)
				if err != nil {
					t.Fatal(err)
				}

				if tt.registerSubject > 1 {
					if err = Register(i, "", binds.Instance(testSecondSubject)); err != nil {
						t.Fatal(err)
					}
				}
				if tt.registerSubject > 0 {
					if err = Register(i, "", binds.Instance(testFirstSubject)); err != nil {
						t.Fatal(err)
					}
				}

				var result string
				result, err = Get[string](i, "")
				if err != nil && !errors.Is(err, tt.expectedError) {
					t.Fatalf(
						"Error is not the same:\nExpected: `%v`\nReceived: `%v`",
						tt.expectedError, err,
					)
				}

				if result != tt.expected {
					t.Error("Result is not the same as expected")
				}
			},
		)
	}
}

func testGuestSubtype[T, K interface{ ~int32 | ~uint8 | ~float64 }](t *testing.T) {
	i := stdinj.New(
		stdinj.Options{Cache: injopts.CacheOptReturnAll, Resolve: injopts.ResolveOptDuckTyping},
	)
	var (
		registerElement K = 0b101010
		expectedElement T // zero value
	)

	if err := Register(i, "", binds.Instance(registerElement)); err != nil {
		t.Fatal(err)
	}

	result, err := Get[T](i, "")
	if err == nil {
		t.Fatalf("No error was received when trying to find subtype `%T`", result)
	}

	if result != expectedElement {
		t.Errorf(
			"Result is not the same as expected\nReceived: `%v`\nExpected: `%v`",
			result, registerElement,
		)
	}
}

func TestGet_guessSubtypes(t *testing.T) {
	type (
		SubTypeInt32   uint8
		SubTypeUint8   uint8
		SubTypeFloat64 float64
	)

	t.Run("Int32 subtype", testGuestSubtype[SubTypeInt32, uint8])
	t.Run("Uint8 subtype", testGuestSubtype[SubTypeUint8, uint8])
	t.Run("Float64 subtype", testGuestSubtype[SubTypeFloat64, uint8])
}

func TestGetAll_withGeneratedBind(t *testing.T) {
	const expectedLanguage = "Portuguese"
	i := stdinj.New(
		stdinj.Options{Cache: injopts.CacheOptReturnAll, Resolve: injopts.ResolveOptDuckTyping},
	)
	err := Register(
		i, "",
		binds.Factory(func(retriever types.DependencyRetriever) (fixtures.CountryLanguage, error) {
			return fixtures.CountryLanguage{Language: expectedLanguage}, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	var result fixtures.Language
	if result, err = Get[fixtures.Language](i, ""); err != nil {
		t.Fatalf("Should not have gotten error when trying to find all subtypes")
	}

	if result.Name() != expectedLanguage {
		t.Errorf(
			"Result is not the same as expected\nReceived: `%v`\nExpected: `%v`",
			result, expectedLanguage,
		)
	}
}

func TestCheckSavedAsBind_pointerTypeDuckTyping(t *testing.T) {
	// This test verifies that checkSavedAsBind correctly handles pointer types
	// when checking against interfaces via duck typing.

	// Register a pointer type using Factory so the bind is stored (not the generated value)
	langPtr := &fixtures.GoProgrammingLang{}
	bind := binds.Factory(
		func(retriever types.DependencyRetriever) (*fixtures.GoProgrammingLang, error) {
			return langPtr, nil
		},
	)

	// Test checkSavedAsBind directly with interface Language
	// This should succeed because PointerValue() fallback allows correct assertion
	result, err := checkSavedAsBind[fixtures.Language](nil, bind)
	if err != nil {
		t.Fatalf("checkSavedAsBind failed with error: %v", err)
	}

	if result == nil {
		t.Fatal("checkSavedAsBind returned nil result, expected valid Language interface")
	}

	// Verify the result is correct
	if (*result).Name() != langPtr.Name() {
		t.Errorf(
			"Language name mismatch. Expected: `%s`, Received: `%s`",
			langPtr.Name(), (*result).Name(),
		)
	}

	if (*result).Kind() != langPtr.Kind() {
		t.Errorf(
			"Language kind mismatch. Expected: `%s`, Received: `%s`",
			langPtr.Kind(), (*result).Kind(),
		)
	}

	interfaceBind := binds.Factory(
		func(retriever types.DependencyRetriever) (fixtures.Language, error) {
			return langPtr, nil
		},
	)
	result, err = checkSavedAsBind[fixtures.Language](nil, interfaceBind)
	if err != nil {
		t.Fatalf("checkSavedAsBind failed with error: %v", err)
	}
	if result != nil {
		t.Fatal("checkSavedAsBind returned no-nil result, expected a nil result")
	}
}

func BenchmarkGet_StoredAsValue(b *testing.B) {
	inj := stdinj.New(stdinj.Options{})
	if err := Register(inj, "", binds.Instance(42)); err != nil {
		b.Fatalf("register: %v", err)
	}

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = Get[int](inj, "")
	}
	runtime.KeepAlive(sink)
}

func BenchmarkGet_StoredAsBind(b *testing.B) {
	inj := stdinj.New(stdinj.Options{})
	if err := Register(inj, "", binds.Factory(
		func(types.DependencyRetriever) (int, error) { return 42, nil },
	)); err != nil {
		b.Fatalf("register: %v", err)
	}

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = Get[int](inj, "")
	}
	runtime.KeepAlive(sink)
}

func BenchmarkGet_InterfaceMissWithoutDuckTyping(b *testing.B) {
	inj := stdinj.New(stdinj.Options{})

	var sink fixtures.Language
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = Get[fixtures.Language](inj, "")
	}
	runtime.KeepAlive(sink)
}
