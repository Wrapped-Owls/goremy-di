package injector

import (
	"errors"
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
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

				i := New(injopts.CacheOptAllowOverride, types.ReflectionOptions{})
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

				i := New(injopts.CacheOptAllowOverride, types.ReflectionOptions{})
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
	// Checks if panics when trying to override
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Function did not panic")
			t.FailNow()
		}
	}()
	inj := New(injopts.CacheOptNone, types.ReflectionOptions{})
	const (
		expectedString   = "avocado"
		unexpectedString = "banana"
	)
	_ = Register(
		inj, binds.Instance(expectedString),
	)

	if result := TryGet[string](inj); result != expectedString {
		t.Error("Instance register is not working as expected")
		t.FailNow()
	}

	_ = Register(
		inj, binds.Singleton(
			func(retriever types.DependencyRetriever) (string, error) {
				return unexpectedString, nil
			},
		),
	)

	if result := TryGet[string](inj); result != expectedString {
		t.Error("Instance bind is being overridden by singleton bind")
	}
}

func TestGetGen(t *testing.T) {
	const expected = "I love Go, yes this is true, as the answer 42"

	interfaceValue := fixtures.GoProgrammingLang{}
	testCases := [...]struct {
		name           string
		getGenCallback func(ij types.Injector) string
		useReflection  bool
	}{
		{
			name:          "GetGen[string]",
			useReflection: true,
			getGenCallback: func(ij types.Injector) string {
				return TryGetGen[string](
					ij,
					[]types.InstancePair[any]{
						{
							Value: uint8(42),
						},
						{
							Value: "Go",
							Key:   "lang",
						},
						{
							Value: true,
						},
						{
							Value:          interfaceValue,
							InterfaceValue: (*fixtures.Language)(nil),
						},
					},
				)
			},
		},
		{
			name: "GetGenFunc[string]",
			getGenCallback: func(i types.Injector) string {
				return TryGetGenFunc[string](
					i, func(ij types.Injector) error {
						err := Register(ij, binds.Instance[uint8](42))
						err = Register(ij, binds.Instance("Go"), "lang")
						err = Register(ij, binds.Instance(true))
						err = Register[fixtures.Language](
							ij,
							binds.Instance[fixtures.Language](interfaceValue),
						)
						return err
					},
				)
			},
		},
	}

	for _, tCase := range testCases {
		i := New(
			injopts.CacheOptAllowOverride,
			types.ReflectionOptions{UseReflectionType: tCase.useReflection},
		)
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

				// Check if the binds doesn't exist after do the GetGen
				var (
					uintResult = TryGet[uint8](i)
					boolResult = TryGet[bool](i)
					strResult  = TryGet[string](i, "lang")
				)
				if uintResult != 0 || boolResult || len(strResult) > 0 {
					t.Error("Parameter injection values override the original injector")
				}
			},
		)
	}
}

func TestGetGen_raiseCastError(t *testing.T) {
	var (
		i = New(
			injopts.CacheOptAllowOverride,
			types.ReflectionOptions{UseReflectionType: true},
		)
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
			_, err = GetGen[string](
				i,
				[]types.InstancePair[any]{
					{
						Value:          interfaceValue,
						InterfaceValue: (*fixtures.Language)(nil),
					},
				},
			)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		},
	)

	t.Run(
		"Register pointer interface value", func(t *testing.T) {
			_, err = GetGen[string](
				i,
				[]types.InstancePair[any]{
					{
						Value:          &interfaceValue,
						InterfaceValue: (*fixtures.Language)(nil),
					},
				},
			)
			if err == nil {
				t.Error("No error has returned after binding the value incorrectly")
				t.FailNow()
			}

			if !errors.Is(err, utils.ErrTypeCastInRuntime) {
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
				expectedError:   utils.ErrElementNotRegistered,
			},
			{
				name:            "Inject multiple elements that implements interface",
				registerSubject: 2,
				expected:        "",
				expectedError:   utils.ErrFoundMoreThanOneValidDI,
			},
		}
	)

	for _, tt := range testCases {
		t.Run(
			tt.name, func(t *testing.T) {
				i := New(injopts.CacheOptReturnAll, types.ReflectionOptions{})
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
	i := New(injopts.CacheOptReturnAll, types.ReflectionOptions{})
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
	i := New(injopts.CacheOptReturnAll, types.ReflectionOptions{})
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
