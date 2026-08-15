package injcontainer

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stdinj"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
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

				i := stdinj.New(stdinj.Options{Cache: injopts.CacheOptAllowOverride})
				if err := Register(i, "", insBind); err != nil {
					t.Error(err)
					t.FailNow()
				}
				for index := 0; index < totalExecutions; index++ {
					result, err := Get[string](i, "")
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

				i := stdinj.New(stdinj.Options{Cache: injopts.CacheOptAllowOverride})
				if invocations != 0 {
					t.Error("Singleton was generated before register")
				}
				for index := 0; index < 11; index++ {
					_ = Register(i, "", sgtBind)
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
					result, err := Get[*string](i, "")
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
	inj := stdinj.New(stdinj.Options{})
	const (
		expectedString   = "avocado"
		unexpectedString = "banana"
	)
	err := Register(
		inj, "", binds.Instance(expectedString),
	)
	if err != nil {
		t.Errorf("Unable to fist register instance: %v", err)
	}

	if result := TryGet[string](inj, ""); result != expectedString {
		t.Error("Instance register is not working as expected")
		t.FailNow()
	}

	err = Register(
		inj, "", binds.Singleton(
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

	if result := TryGet[string](inj, ""); result != expectedString {
		t.Error("Instance bind is being overridden by singleton bind")
	}
}
