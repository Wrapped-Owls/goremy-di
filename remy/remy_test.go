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
