package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"math"
	"testing"
)

func TestInjection__GetNoRegistered(t *testing.T) {
	ij := New(false, types.ReflectionOptions{})

	if strResult := Get[string](ij); len(strResult) != 0 {
		t.Errorf("string result is not the default, received: `%s`", strResult)
	}
	if intResult := Get[int](ij); intResult != 0 {
		t.Errorf("int result is not the default, received: %d", intResult)
	}
	if pointerResult := Get[*bool](ij); pointerResult != nil {
		t.Error("pointer received is not null")
	}
	if interfaceResult := Get[interface{ a() string }](ij); interfaceResult != nil {
		t.Error("interface result is not nil")
	}
	if structResult := Get[struct{ element string }](ij); len(structResult.element) != 0 {
		t.Error("default struct is not created correctly")
	}
}

type guide struct {
	value string
}

func (g guide) String() string {
	return g.value
}

func TestInjection__GetStructImplementInterface(t *testing.T) {
	expected := [...]guide{
		{value: "DON'T PANIC"},
		{value: "DO PANIC"},
	}
	type universalAnswer interface {
		String() string
	}
	ij := New(false, types.ReflectionOptions{GenerifyInterface: true})

	Register(ij, binds.Instance(
		func(retriever types.DependencyRetriever) universalAnswer {
			return &expected[0]
		},
	))
	// Register again as another type, to check if it works
	Register(ij, binds.Instance(
		func(retriever types.DependencyRetriever) guide {
			return expected[1]
		},
	))

	result := Get[universalAnswer](ij)
	if result != &expected[0] {
		t.Errorf("element injected is different than the provided. Received %p", result)
	} else if result.String() != expected[0].value {
		t.Errorf("element was reseted. Expected: `%s`; Received: `%s`", expected[0].value, result.String())
	}

	structResult := Get[guide](ij)
	if structResult.String() != expected[1].value {
		t.Errorf("element was reseted. Expected: `%s`; Received: `%s`", expected[1].value, structResult.String())
	}
}

func TestInjection__RegisterSameKeyDifferentType(t *testing.T) {
	const (
		expectedStr = "DON'T PANIC"
		expectedInt = 42
	)

	ij := New(false, types.ReflectionOptions{})
	Register(
		ij,
		binds.Instance(func(retriever types.DependencyRetriever) string {
			return expectedStr
		}),
		"truth",
	)
	Register(
		ij,
		binds.Instance(func(retriever types.DependencyRetriever) int {
			return expectedInt
		}),
		"truth",
	)

	strResult := Get[string](ij, "truth")
	intResult := Get[int](ij, "truth")

	if strResult != expectedStr {
		t.Errorf("string injection should not be overrided. Received: `%s`. Expected: `%s`", strResult, expectedStr)
	}
	if intResult != expectedInt {
		t.Errorf("int injection should not be overrided. Received: `%d`. Expected: `%d`", intResult, expectedInt)
	}

}

func TestInjection__RetrieveSameTypeDifferentKey(t *testing.T) {
	var (
		resultParts = [...]string{
			"I'm programming in ",
			"go",
		}
	)
	a := binds.Instance(
		func(ij types.DependencyRetriever) string {
			language := Get[string](ij, "lang")
			return resultParts[0] + language
		},
	)

	ij := New(true, types.ReflectionOptions{})
	Register(
		ij,
		binds.Instance(func(retriever types.DependencyRetriever) string {
			return resultParts[1]
		}),
		"lang",
	)
	Register(ij, a)
	result := Get[string](ij)

	if result != resultParts[0]+resultParts[1] {
		t.Errorf("injection result not work properly: Received: `%s`", result)
	}
}

type speakerImpl struct {
	sound string
}

func (s speakerImpl) speak() string {
	return s.sound
}

func TestInjection__RegisterEqualInterfaces(t *testing.T) {
	type (
		spk1 interface {
			speak() string
		}
		spk2 interface {
			spk1
		}
	)

	const storageKey = "same-key"
	elements := [...]speakerImpl{
		{sound: "meow"},
		{sound: "woof woof"},
	}

	ij := New(true, types.ReflectionOptions{})
	Register(
		ij,
		binds.Instance(func(retriever types.DependencyRetriever) spk1 {
			return elements[0]
		}),
		storageKey,
	)
	Register(
		ij,
		binds.Instance(func(retriever types.DependencyRetriever) spk2 {
			return elements[1]
		}),
		storageKey,
	)

	// Start to retrieve the injected objects
	results := [...]spk1{
		Get[spk1](ij, storageKey),
		Get[spk2](ij, storageKey),
	}

	if results[0] == nil || results[1] == nil {
		t.Error("Failed to retrieve injected objects thought interfaces")
		t.FailNow()
	}
	for index := 0; index < int(math.Min(float64(len(results)), float64(len(elements)))); index++ {
		if results[index].speak() != elements[index].speak() {
			t.Errorf(
				"element got by interface is not performing correctly: [%d] `%s` != `%s`",
				index, results[index].speak(), elements[index].speak(),
			)
		}
	}
}
