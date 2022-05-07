package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/types"
	"testing"
)

// TestGenerateBind__InstanceFactory verify if when registering an instance, it is only generated once
func TestGenerateBind__InstanceFactory(t *testing.T) {
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
			name:                "instance",
			expectedGenerations: 1,
			bindGenerator:       binds.Instance[string],
		},
		{
			name:                "factory",
			expectedGenerations: totalExecutions,
			bindGenerator:       binds.Factory[string],
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			counter := 0
			insBind := c.bindGenerator(func(retriever types.DependencyRetriever) string {
				counter++
				return expectedString
			})

			i := New()
			Register[string](i, insBind)
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

func TestRegister__Singleton(t *testing.T) {
	const totalGetsExecuted = 11

	cases := []struct {
		name                string
		expected            string
		registerGenerations int
		bindGenerator       func(types.Binder[*string]) types.Bind[*string]
	}{
		{
			name:                "singleton",
			expected:            "here we go",
			registerGenerations: 1,
			bindGenerator:       binds.Singleton[*string],
		},
		{
			name:                "lazy singleton",
			expected:            "JUST BE SURE TO LAZY",
			registerGenerations: 0,
			bindGenerator:       binds.LazySingleton[*string],
		},
	}

	for _, bindCase := range cases {
		t.Run(bindCase.expected, func(t *testing.T) {
			var (
				invocations = 0
			)
			sgtBind := bindCase.bindGenerator(func(retriever types.DependencyRetriever) *string {
				invocations++
				return &bindCase.expected
			})

			i := New()
			if invocations != 0 {
				t.Error("Singleton was generated before register")
			}
			for index := 0; index < 11; index++ {
				Register[*string](i, sgtBind)
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
