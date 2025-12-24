package injector

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

// TestGenerateBind__InstanceFactory verify if when registering an instance, it is only generated once
func TestGenerateBind__InstanceFactory(testObj *testing.T) {
	const (
		expectedString  = "avocado"
		totalExecutions = 11
	)

	cases := []struct {
		name                string
		expectedGenerations int
		bindGenerator       func(func() string) types.Bind[string]
	}{
		{
			name:                "INSTANCE",
			expectedGenerations: 1,
			bindGenerator: func(factory func() string) types.Bind[string] {
				return binds.Instance[string](factory())
			},
		},
		{
			name:                "FACTORY",
			expectedGenerations: totalExecutions,
			bindGenerator: func(factory func() string) types.Bind[string] {
				return binds.Factory[string](
					func(retriever types.DependencyRetriever) (string, error) {
						return factory(), nil
					},
				)
			},
		},
	}

	for _, c := range cases {
		testObj.Run(
			c.name, func(t *testing.T) {
				counter := 0
				insBind := c.bindGenerator(
					func() string {
						counter++
						return expectedString
					},
				)

				i := New(injopts.CacheOptAllowOverride)
				if err := Register(i, insBind); err != nil {
					t.Error(err)
					t.FailNow()
				}
				for index := 0; index < totalExecutions; index++ {
					result, err := Get[string](i)
					if result != expectedString {
						t.Error("Generated instance is incorrect")
					}
					if err != nil {
						t.Error(err)
					}
				}

				if counter != c.expectedGenerations {
					t.Errorf("Bind generated %d times. Expected %d", counter, c.expectedGenerations)
				}
			},
		)
	}
}

func TestRegister__Singleton(testObj *testing.T) {
	const totalGetsExecuted = 11

	cases := []struct {
		name                string
		expected            string
		registerGenerations int
		bindGenerator       func(types.Binder[*string]) types.Bind[*string]
	}{
		{
			name:                "SINGLETON",
			expected:            "here we go",
			registerGenerations: 1,
			bindGenerator:       binds.Singleton[*string],
		},
		{
			name:                "LAZY_SINGLETON",
			expected:            "JUST BE SURE TO LAZY",
			registerGenerations: 0,
			bindGenerator:       binds.LazySingleton[*string],
		},
	}

	for _, bindCase := range cases {
		testObj.Run(
			bindCase.name, func(t *testing.T) {
				invocations := 0
				sgtBind := bindCase.bindGenerator(
					func(retriever types.DependencyRetriever) (*string, error) {
						invocations++
						return &bindCase.expected, nil
					},
				)

				i := New(injopts.CacheOptAllowOverride)
				if invocations != 0 {
					t.Error("Singleton was generated before register")
				}
				for index := 0; index < 11; index++ {
					_ = Register(i, sgtBind)
					if invocations != bindCase.registerGenerations {
						t.Errorf(
							"Singleton %d times. Expected %d",
							invocations,
							bindCase.registerGenerations,
						)
						t.FailNow()
					}
				}

				for index := 0; index < totalGetsExecuted; index++ {
					result, err := Get[*string](i)
					if err != nil {
						t.Error(err)
					}
					if result != &bindCase.expected {
						t.Errorf("Singleton is not working as singleton")
					}
					if invocations != 1 {
						t.Errorf("Singleton generated %d times", invocations)
					}
				}
			},
		)
	}
}

// TestRegister__overrideInstanceByBind verify if when overriding a instance
func TestRegister__overrideInstanceByBind(t *testing.T) {
	inj := New(injopts.CacheOptNone)
	const (
		expectedString   = "avocado"
		unexpectedString = "banana"
	)
	err := Register(
		inj, binds.Instance(expectedString),
	)
	if err != nil {
		t.Errorf("Unable to fist register instance: %v", err)
	}

	if result := TryGet[string](inj); result != expectedString {
		t.Error("Instance register is not working as expected")
		t.FailNow()
	}

	err = Register(
		inj, binds.Singleton(
			func(retriever types.DependencyRetriever) (string, error) {
				return unexpectedString, nil
			},
		),
	)
	if err == nil {
		t.Fatalf("Instance was registered unexpectedly")
	} else if !errors.Is(err, remyErrs.ErrAlreadyBoundSentinel) {
		t.Errorf("Result error is not the expected error: %v", err.Error())
	}

	if result := TryGet[string](inj); result != expectedString {
		t.Error("Instance bind is being overridden by singleton bind")
	}
}

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
					ij,
					[]types.BindEntry{
						types.NewBindPair(uint8(42), ""),
						types.NewBindPair("Go", "lang"),
						types.NewBindPair(true, ""),
						types.NewBindPair[fixtures.Language](interfaceValue, ""),
					},
				)
				return result
			},
		},
		{
			name: "GetWith[string]",
			getGenCallback: func(i types.Injector) string {
				result, _ := GetWith[string](
					i, func(ij types.Injector) error {
						err := errors.Join(
							Register(ij, binds.Instance[uint8](42)),
							Register(ij, binds.Instance("Go"), "lang"),
							Register(ij, binds.Instance(true)),
							Register[fixtures.Language](
								ij, binds.Instance[fixtures.Language](interfaceValue),
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
		i := New(injopts.CacheOptAllowOverride)
		_ = Register(
			i, binds.Factory(
				func(retriever types.DependencyRetriever) (result string, err error) {
					result = fmt.Sprintf(
						"I love %s, yes this is %v, as the answer %d",
						TryGet[string](
							retriever,
							"lang",
						),
						TryGet[bool](retriever),
						TryGet[uint8](retriever),
					)

					if _, err = Get[fixtures.Language](retriever); err != nil {
						t.Error(err)
					}
					return
				},
			),
		)

		// register a bool bind to check if it will be replaced during parameter passing
		_ = Register(i, binds.Instance(false))

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

				// Check if the binds doesn't exist after do the GetWithPairs
				var (
					uintResult, _ = Get[uint8](i)
					boolResult, _ = Get[bool](i)
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

	i := New(injopts.CacheOptAllowOverride)

	errFirstRegister := errors.Join(
		Register(
			i, binds.Factory(
				func(retriever types.DependencyRetriever) (result string, err error) {
					workTime := TryGet[time.Time](retriever)
					timeStr := workTime.Format("3 PM")
					result = fmt.Sprintf(
						"%s and %s work at the park at %s, during: %d minutes, is weekend: %v",
						TryGet[string](
							retriever,
							"employee1",
						),
						TryGet[string](retriever, "employee2"),
						timeStr,
						TryGet[uint8](retriever),
						TryGet[bool](retriever),
					)
					return
				},
			),
		),

		// register a bool bind to check if it will be replaced during parameter passing
		Register(i, binds.Instance(false)),
		Register(i, binds.Instance(time.Time{})),
	)
	if errFirstRegister != nil {
		t.Fatal(errFirstRegister)
	}

	// Test with direct BindKey provided - when Key is provided, InterfaceValue is not needed
	result, err := GetWithPairs[string](
		i, []types.BindEntry{
			types.NewBindPair(uint8(42), ""),
			types.NewBindPair("Mordecai", "employee1"),
			types.NewBindPair("Rigby", "employee2"),
			types.NewBindPair(time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC), ""),
			types.NewBindPair(true, ""),
		},
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
		uintResult, _      = Get[uint8](i)
		boolResult, _      = Get[bool](i)
		employee1Result, _ = Get[string](i, "employee1")
		employee2Result, _ = Get[string](i, "employee2")
		timeResult, _      = Get[time.Time](i)
	)
	if uintResult != 0 || boolResult || len(employee1Result) > 0 ||
		len(employee2Result) > 0 || !timeResult.IsZero() {
		t.Error("Parameter injection values override the original injector")
	}
}

func TestGetGen_raiseCastError(t *testing.T) {
	var (
		i                                = New(injopts.CacheOptAllowOverride)
		interfaceValue fixtures.Language = fixtures.GoProgrammingLang{}
	)
	err := Register(
		i, binds.Factory(
			func(retriever types.DependencyRetriever) (result string, getErr error) {
				var lang fixtures.Language
				if lang, getErr = Get[fixtures.Language](retriever); getErr == nil {
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
				i, []types.BindEntry{types.NewBindPair[fixtures.Language](interfaceValue, "")},
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
				i,
				[]types.BindEntry{
					types.InstancePair[*fixtures.Language]{
						Key:   utils.NewKeyElem[fixtures.Language](),
						Value: &interfaceValue,
					},
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
				i := New(injopts.CacheOptReturnAll)
				err := Register(
					i, binds.Factory(
						func(retriever types.DependencyRetriever) (result string, getErr error) {
							var lang fixtures.Language
							if lang, getErr = Get[fixtures.Language](retriever); getErr == nil {
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
					if err = Register(i, binds.Instance(testSecondSubject)); err != nil {
						t.Fatal(err)
					}
				}
				if tt.registerSubject > 0 {
					if err = Register(i, binds.Instance(testFirstSubject)); err != nil {
						t.Fatal(err)
					}
				}

				var result string
				result, err = Get[string](i)
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
	i := New(injopts.CacheOptReturnAll)
	var (
		registerElement K = 0b101010
		expectedElement T // zero value
	)

	if err := Register(i, binds.Instance(registerElement)); err != nil {
		t.Fatal(err)
	}

	result, err := Get[T](i)
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
	i := New(injopts.CacheOptReturnAll)
	err := Register(
		i,
		binds.Factory(func(retriever types.DependencyRetriever) (fixtures.CountryLanguage, error) {
			return fixtures.CountryLanguage{Language: expectedLanguage}, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	var result fixtures.Language
	if result, err = Get[fixtures.Language](i); err != nil {
		t.Fatalf("Should not have gotten error when trying to find all subtypes")
	}

	if result.Name() != expectedLanguage {
		t.Errorf(
			"Result is not the same as expected\nReceived: `%v`\nExpected: `%v`",
			result, expectedLanguage,
		)
	}
}

func TestGetWith_withParentDuckTyping(t *testing.T) {
	// Create parent injector with CacheOptReturnAll enabled (allows GetAll)
	parent := New(injopts.CacheOptReturnAll)

	// Register an interface implementation in the parent injector
	langImpl := fixtures.GoProgrammingLang{}
	if err := Register(parent, binds.Instance(langImpl)); err != nil {
		t.Fatalf("Failed to register language implementation: %v", err)
	}

	// Use GetWith which creates a sub-injector with CacheOptNone (doesn't allow GetAll)
	// The sub-injector should be able to find the interface via duck typing by delegating to parent
	result, err := GetWith[fixtures.Language](
		parent, func(ij types.Injector) error {
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
