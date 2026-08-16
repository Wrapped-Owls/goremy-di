package binds

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type elemViews interface {
	PointerValue() any
	DefaultValue() any
	ElementKey() types.BindKey
}

func TestElemType_Views(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind types.Bind[string]
	}{
		{name: "through a factory bind", bind: Factory(constantBinder(""))},
		{name: "through an instance bind", bind: Instance("")},
		{name: "through a singleton bind", bind: Singleton(constantBinder(""))},
		{name: "through a lazy singleton bind", bind: LazySingleton(constantBinder(""))},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			views, ok := testCase.bind.(elemViews)
			if !ok {
				t.Fatalf("%T does not carry the elemType views", testCase.bind)
			}

			pointer, isPointer := views.PointerValue().(*string)
			if !isPointer || pointer != nil {
				t.Fatalf("PointerValue() = %#v, want a nil *string", views.PointerValue())
			}
			if got := views.DefaultValue(); got != "" {
				t.Fatalf("DefaultValue() = %#v, want the zero string", got)
			}
			if got := views.ElementKey(); got != (types.KeyElem[string]{}) {
				t.Fatalf("ElementKey() = %#v, want KeyElem[string]", got)
			}
		})
	}
}

// the views must describe T itself, never the value a bind happens to hold
func TestElemType_ViewsFollowTheTypeParameter(t *testing.T) {
	t.Parallel()

	views, ok := Instance(42).(elemViews)
	if !ok {
		t.Fatal("FactoryBind[int] does not carry the elemType views")
	}

	if got := views.DefaultValue(); got != 0 {
		t.Fatalf("DefaultValue() = %#v, want the zero int even though the bind holds 42", got)
	}
	if got := views.ElementKey(); got == (types.KeyElem[string]{}) {
		t.Fatal("ElementKey() for int must not equal the string key")
	}
}
