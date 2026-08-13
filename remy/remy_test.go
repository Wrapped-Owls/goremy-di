package remy

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestOverride(t *testing.T) {
	var checkpoints uint8 = 0
	// Checks if panics when trying to override
	defer func() {
		r := recover()
		if checkpoints == 0 {
			if r == nil {
				t.Error("Function did not panic")
				t.FailNow()
			}
		}
		if checkpoints != 1 {
			t.Error("Test panic on wrong checkpoint")
		}
	}()

	// create an injector that can override a bind and try to register it twice
	inj := NewInjector(Config{CanOverride: true})
	RegisterInstance(inj, "test")
	Override(inj, Instance("test_override"))
	checkpoints++
	RegisterInstance(inj, "test_panic_override")
	checkpoints++
}

// TestOverride__panicIfNotAllowed executes a test to check the rule that is:
// "When Override is not allowed in the injector, the function should panic when trying to override by any method"
func TestOverride__panicIfNotAllowed(t *testing.T) {
	var checkpoints uint8 = 0
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Function did not panic")
			t.FailNow()
		}
		if checkpoints != 0 {
			t.Error("Test panic after reaching the first checkpoint")
		}
	}()

	inj := NewInjector(Config{CanOverride: false})
	RegisterInstance(inj, "test")
	Override(inj, Instance("test_override"))
	checkpoints++
}

func TestRegisterGet(t *testing.T) {
	i := NewInjector(Config{DuckTypeElements: true})
	t.Run("Int bind with GetAll", func(t *testing.T) {
		RegisterInstance[int](i, 42)
		result := MustGetAll[int](i)
		if len(result) != 1 || result[0] != 42 {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("String bind with Get", func(t *testing.T) {
		RegisterInstance(i, "hello", "greeting")
		result := MustGet[string](i, "greeting")
		if result != "hello" {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

func TestRegisterSingleton(t *testing.T) {
	i := NewInjector(Config{DuckTypeElements: true})
	var totalCalls uint16
	// It Should run only once during register and after it, the call only returns generated value
	RegisterSingleton(i, func(retriever DependencyRetriever) (uint16, error) {
		totalCalls += 1
		return totalCalls, nil
	})

	if totalCalls != 1 {
		t.Errorf("Expected total calls to be 1, but got %d", totalCalls)
	}

	// Get the value multiple times and verify that it's always 1
	for index := 0; index < 10; index++ {
		value := MustGet[uint16](i)
		if value != 1 {
			t.Errorf("Expected value to be 1, but got %d", value)
		}
	}
}

func TestGetWithContext_ReturnsValue(t *testing.T) {
	type testKey struct{}
	ctxKey := testKey{}
	inj := NewInjector()

	RegisterConstructorArgs1(inj, Factory[string], func(ctx context.Context) string {
		val, _ := ctx.Value(ctxKey).(string)
		return val
	})

	const injectedValue = "the-Blade+!"
	ctx := context.WithValue(context.Background(), ctxKey, injectedValue)

	result, err := GetWithContext[string](inj, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != injectedValue {
		t.Fatalf("unexpected value: %s", result)
	}
}

func TestGetWithContext_Error(t *testing.T) {
	inj := NewInjector()

	_, err := GetWithContext[string](inj, context.Background())
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !errors.Is(err, ErrElementNotRegistered) {
		t.Fatalf("expected ErrElementNotRegistered but got %v", err)
	}
}

type multiBindGreeter interface{ Greet() string }

type multiBindEnglish struct{}

func (multiBindEnglish) Greet() string { return "hello" }

func TestConfig_MultiBinding(t *testing.T) {
	testCases := []struct {
		name          string
		config        Config
		wantGetAllErr error
		wantDuckTyped bool
	}{
		{
			name:          "multi binding enables GetAll without implicit discovery",
			config:        Config{MultiBinding: true},
			wantGetAllErr: nil,
			wantDuckTyped: false,
		},
		{
			name:          "duck typing implies multi binding and implicit discovery",
			config:        Config{DuckTypeElements: true},
			wantGetAllErr: nil,
			wantDuckTyped: true,
		},
		{
			name:   "default config allows neither",
			config: Config{},
			wantGetAllErr: errors.New(
				"the current injector config does not allow returning all elements",
			),
			wantDuckTyped: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			inj := NewInjector(testCase.config)
			Register(inj, Instance(multiBindEnglish{}))
			RegisterInstance(inj, "first", "one")
			RegisterInstance(inj, "second", "two")

			_, getAllErr := GetAll[string](inj, "one")
			if (getAllErr != nil) != (testCase.wantGetAllErr != nil) {
				t.Fatalf(
					"GetAll error = %v, want error presence %v",
					getAllErr,
					testCase.wantGetAllErr != nil,
				)
			}

			_, duckErr := Get[multiBindGreeter](inj)
			if duckTyped := duckErr == nil; duckTyped != testCase.wantDuckTyped {
				t.Fatalf(
					"implicit interface discovery = %v (err %v), want %v",
					duckTyped,
					duckErr,
					testCase.wantDuckTyped,
				)
			}
		})
	}
}

func TestGet_DependencyTrace(t *testing.T) {
	inj := NewInjector()

	// Register string that depends on int
	Register(inj, Factory(func(r DependencyRetriever) (string, error) {
		_, err := Get[int](r)
		return "", err
	}))

	_, err := Get[string](inj)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errMsg := err.Error()
	// The message should contain both "string" and "int" as they are in the resolution path
	if !strings.Contains(errMsg, "string") || !strings.Contains(errMsg, "int") {
		t.Errorf("error message does not contain all types in the path: %s", errMsg)
	}

	// Check if we can still identify the root cause using errors.Is
	if !errors.Is(err, ErrElementNotRegistered) {
		t.Errorf("root cause ErrElementNotRegistered lost: %v", err)
	}

	// Test with deeper dependencies: string -> int -> bool
	Register(inj, Factory(func(r DependencyRetriever) (int, error) {
		_, err := Get[bool](r)
		return 0, err
	}))

	_, err = Get[string](inj)
	errMsg = err.Error()
	if !strings.Contains(errMsg, "string") || !strings.Contains(errMsg, "int") ||
		!strings.Contains(errMsg, "bool") {
		t.Errorf("deep error message does not contain all types in the path: %s", errMsg)
	}
}
