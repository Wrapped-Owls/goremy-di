package injector

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/types"
	"testing"
)

// TestRegister__Instance verify if when registering an instance, it is only generated once
func TestGenerateBind__Instance(t *testing.T) {
	const expected = "avocado"
	counter := 0
	insBind := binds.Instance[string](func(retriever types.DependencyRetriever) string {
		counter++
		return expected
	})

	i := New()
	Register[string](i, insBind)
	for index := 0; index < 11; index++ {
		result := Get[string](i)
		if result != expected {
			t.Error("Generated instance is incorrect")
		}
	}
	if counter > 1 {
		t.Errorf("Instance bind generated %d times. Expected 1", counter)
	}
}

func TestGenerateBind__Factory(t *testing.T) {
	const (
		expectedString      = "avocado"
		expectedGenerations = 11
	)
	counter := 0
	insBind := binds.Factory[string](func(retriever types.DependencyRetriever) string {
		counter++
		return expectedString
	})

	i := New()
	Register[string](i, insBind)
	for index := 0; index < expectedGenerations; index++ {
		result := Get[string](i)
		if result != expectedString {
			t.Error("Generated instance is incorrect")
		}
	}
	if counter != expectedGenerations {
		t.Errorf("Instance bind generated %d times. Expected %d", counter, expectedGenerations)
	}
}

func TestRegister__Singleton(t *testing.T) {
	const totalGetsExecuted = 11

	var (
		invocations = 0
		instance    = "testing singleton"
	)
	sgtBind := binds.Singleton[*string](func(retriever types.DependencyRetriever) *string {
		invocations++
		return &instance
	})

	i := New()
	if invocations != 0 {
		t.Error("Singleton was generated before register")
	}
	for index := 0; index < 11; index++ {
		Register[*string](i, sgtBind)
		if invocations != 1 {
			t.Error("Singleton generated more than once")
			t.FailNow()
		}
	}

	for index := 0; index < totalGetsExecuted; index++ {
		result := Get[*string](i)
		if result != &instance {
			t.Errorf("Singleton is not working as singleton")
		}
		if invocations != 1 {
			t.Errorf("Singleton generated %d times", invocations)
		}
	}
}

func TestRegister__LazySingleton(t *testing.T) {
	var (
		invocations = 0
		instance    = "lazy singleton"
	)
	sgtBind := binds.LazySingleton[*string](func(retriever types.DependencyRetriever) *string {
		invocations++
		return &instance
	})

	i := New()
	if invocations != 0 {
		t.Error("Singleton was generated before register")
	}
	for index := 0; index < 11; index++ {
		Register[*string](i, sgtBind)
		if invocations != 0 {
			t.Errorf("Singleton was generated when should not. Received %d, Expected %d", invocations, 0)
			t.FailNow()
		}
	}

	result := Get[*string](i)
	if result != &instance {
		t.Errorf("Singleton is not working as singleton")
	} else if invocations != 1 {
		t.Errorf("Singleton generated %d times", invocations)
	}
}
