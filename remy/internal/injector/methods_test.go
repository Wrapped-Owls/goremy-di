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
		bindGenerator       func(types.Binder[string]) types.Bind[string]
	}{
		{
			name:                "INSTANCE",
			expectedGenerations: 1,
			bindGenerator:       binds.Instance[string],
		},
		{
			name:                "FACTORY",
			expectedGenerations: totalExecutions,
			bindGenerator:       binds.Factory[string],
		},
	}

	for _, c := range cases {
		testObj.Run(c.name, func(t *testing.T) {
			counter := 0
			insBind := c.bindGenerator(func(retriever types.DependencyRetriever) string {
				counter++
				return expectedString
			})

			i := New(true, false)
			Register(i, insBind)
			for index := 0; index < totalExecutions; index++ {
				result := Get[string](i)
				if result != expectedString {
					t.Error("Generated instance is incorrect")
				}
			}

			if counter != c.expectedGenerations {
				t.Errorf("Instance bind generated %d times. Expected %d", counter, c.expectedGenerations)
			}
		})
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
		testObj.Run(bindCase.name, func(t *testing.T) {
			var (
				invocations = 0
			)
			sgtBind := bindCase.bindGenerator(func(retriever types.DependencyRetriever) *string {
				invocations++
				return &bindCase.expected
			})

			i := New(true, false)
			if invocations != 0 {
				t.Error("Singleton was generated before register")
			}
			for index := 0; index < 11; index++ {
				Register(i, sgtBind)
				if invocations != bindCase.registerGenerations {
					t.Errorf("Singleton %d times. Expected %d", invocations, bindCase.registerGenerations)
					t.FailNow()
				}
			}

			for index := 0; index < totalGetsExecuted; index++ {
				result := Get[*string](i)
				if result != &bindCase.expected {
					t.Errorf("Singleton is not working as singleton")
				}
				if invocations != 1 {
					t.Errorf("Singleton generated %d times", invocations)
				}
			}
		})
	}
}

func TestGetGen(t *testing.T) {
	const expected = "I love Go, yes this is true, as the answer 42"

	i := New(true, false)
	Register(
		i, binds.Factory(func(ij types.DependencyRetriever) string {
			return fmt.Sprintf(
				"I love %s, yes this is %v, as the answer %d",
				Get[string](ij, "lang"), Get[bool](ij), Get[uint8](ij),
			)
		}),
	)

	// register a bool bind to check if it will be replaced during parameter passing
	Register(
		i, binds.Instance(func(ij types.DependencyRetriever) bool {
			return false
		}),
	)

	result := GetGen[string](
		i,
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

	if result != expected {
		t.Errorf(
			"The direct params was not injected correctly.\nExpected: `%s`\nReceived: `%s`",
			expected, result,
		)
		t.FailNow()
	}

	// Check if the binds doesn't exist after do the GetGen
	uintResult := Get[uint8](i)
	boolResult := Get[bool](i)
	strResult := Get[string](i, "lang")

	if uintResult != 0 || boolResult || len(strResult) > 0 {
		t.Error("Parameter injection values override the original injector")
	}
}
