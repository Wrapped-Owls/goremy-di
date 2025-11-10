package remy

import "testing"

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
