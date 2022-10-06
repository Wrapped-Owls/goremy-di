package injector

import (
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
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

				i := New(true, types.ReflectionOptions{})
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

				i := New(true, types.ReflectionOptions{})
				if invocations != 0 {
					t.Error("Singleton was generated before register")
				}
				for index := 0; index < 11; index++ {
					_ = Register(i, sgtBind)
					if invocations != bindCase.registerGenerations {
						t.Errorf("Singleton %d times. Expected %d", invocations, bindCase.registerGenerations)
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
	inj := New(false, types.ReflectionOptions{})
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

	testCases := [...]struct {
		name           string
		getGenCallback func(ij types.Injector) string
	}{
		{
			name: "GetGen[string]",
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
						return err
					},
				)
			},
		},
	}

	for _, tCase := range testCases {
		i := New(true, types.ReflectionOptions{})
		_ = Register(
			i, binds.Factory(
				func(ij types.DependencyRetriever) (result string, err error) {
					result = fmt.Sprintf(
						"I love %s, yes this is %v, as the answer %d",
						TryGet[string](ij, "lang"), TryGet[bool](ij), TryGet[uint8](ij),
					)
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
						expected, result,
					)
					t.FailNow()
				}

				// Check if the binds doesn't exist after do the GetGen
				uintResult := TryGet[uint8](i)
				boolResult := TryGet[bool](i)
				strResult := TryGet[string](i, "lang")

				if uintResult != 0 || boolResult || len(strResult) > 0 {
					t.Error("Parameter injection values override the original injector")
				}
			},
		)
	}
}
