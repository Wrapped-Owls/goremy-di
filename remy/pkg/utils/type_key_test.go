package utils

import (
	"testing"
	"unsafe"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	aTypes "github.com/wrapped-owls/goremy-di/remy/test/fixtures/a/testtypes"
	bTypes "github.com/wrapped-owls/goremy-di/remy/test/fixtures/b/testtypes"
)

func TestNewKeyElem__Generify(t *testing.T) {
	type (
		super interface {
			a() bool
			b() string
			c(int) float32
			d(string) struct{ name string }
		}
		sub interface {
			super
		}
	)

	keys := struct{ super, sub types.BindKey }{
		sub:   NewKeyElem[sub](),
		super: NewKeyElem[super](),
	}
	if keys.super == keys.sub {
		t.Error("type names was the same when should not generify")
	}
}

func TestNewKeyElem_GenerifyWithSameDeclaration(t *testing.T) {
	type (
		super interface{ Do() string }
		sub   interface{ Do() string }
		concr struct{ value string } //nolint:unused // only used to test concrete type
	)

	keys := struct{ super, sub, concr types.BindKey }{
		sub:   NewKeyElem[sub](),
		super: NewKeyElem[super](),
		concr: NewKeyElem[concr](),
	}
	if keys.super == keys.sub {
		t.Fatalf("expected interfaces to have different keys (generification removed)")
	}

	if keys.sub == keys.concr {
		t.Fatalf("expected concrete type to be different from interface type")
	}
}

func TestNewKeyElem__SameStructWithDifferentPackage(t *testing.T) {
	// types from different packages will have different keys
	keys := struct{ a, b types.BindKey }{
		a: NewKeyElem[aTypes.Syringe](),
		b: NewKeyElem[bTypes.Syringe](),
	}
	if keys.a == keys.b {
		t.Error(
			"type names was the same, when it should be different, because of different packages",
		)
	}
}

func TestNewKeyElem__Functions(t *testing.T) {
	type (
		voidCallback        = func()
		stringCallback      = func() string
		multiArgsCallback   = func(...uint8) any
		boolCheckerCallback = func(...string) bool

		// Named args
		namedStringCallback      = func() (result string)
		namedMultiArgsCallback   = func(args ...uint8) (result any)
		namedBoolCheckerCallback = func(languages ...string) bool
	)

	// Check for named and unnamed functions
	t.Run(
		"NamedUnnamed functions", func(t *testing.T) {
			cases := [...][2]types.BindKey{
				{NewKeyElem[namedStringCallback](), NewKeyElem[stringCallback]()},
				{NewKeyElem[namedMultiArgsCallback](), NewKeyElem[multiArgsCallback]()},
				{
					NewKeyElem[namedBoolCheckerCallback](),
					NewKeyElem[boolCheckerCallback](),
				},
			}
			for _, results := range cases {
				if results[0] != results[1] {
					t.Errorf(
						"Named and unnamed functions have been identified as different\nExpected: `%s`\nReceived: `%s`",
						results[0],
						results[1],
					)
				}
			}
		},
	)

	// Check for function types that are different
	t.Run(
		"Different function types", func(t *testing.T) {
			cases := [...][2]types.BindKey{
				{NewKeyElem[namedStringCallback](), NewKeyElem[voidCallback]()},
				{NewKeyElem[stringCallback](), NewKeyElem[multiArgsCallback]()},
				{NewKeyElem[voidCallback](), NewKeyElem[boolCheckerCallback]()},
			}

			for _, results := range cases {
				if results[0] == results[1] {
					t.Errorf(
						"Function types should be different\nFunc_1: `%s`\nFunc_2: `%s`",
						results[0], results[1],
					)
				}
			}
		},
	)

	// Check for function pointers
	var keys struct{ raw, ptr types.BindKey }
	keys.raw = NewKeyElem[namedStringCallback]()
	keys.ptr = NewKeyElem[*namedStringCallback]()
	if keys.raw == keys.ptr {
		t.Error("Function pointer should be different than function type")
	}
}

func TestIsInterface(t *testing.T) {
	type myStruct struct{}
	type myInterface interface{ Foo() }
	type embedded interface {
		error
	}

	// pointer alias
	type ptrInt = *int

	tests := []struct {
		name     string
		fn       any
		expected bool
	}{
		// primitive
		{"int", IsInterface[int], false},
		{"string", IsInterface[string], false},
		{"bool", IsInterface[bool], false},
		{"float64", IsInterface[float64], false},
		{"uintptr", IsInterface[uintptr], false},

		// struct + named struct
		{"struct{}", IsInterface[struct{}], false},
		{"named struct", IsInterface[myStruct], false},

		// pointer types
		{"*int", IsInterface[*int], false},
		{"pointer alias", IsInterface[ptrInt], false},
		{"pointer to struct", IsInterface[*myStruct], false},
		{"pointer to interface type", IsInterface[*error], false},
		{"pointer to slice", IsInterface[*[]int], false},
		{"unsafe.Pointer", IsInterface[unsafe.Pointer], false},

		// slice / array / map
		{"slice", IsInterface[[]int], false},
		{"array", IsInterface[[3]int], false},
		{"map", IsInterface[map[string]int], false},

		// chan
		{"chan", IsInterface[chan int], false},
		{"receive-only chan", IsInterface[<-chan int], false},
		{"send-only chan", IsInterface[chan<- int], false},

		// funcs
		{"func()", IsInterface[func()], false},
		{"func returning int", IsInterface[func() int], false},
		{"nilable func", IsInterface[func(int) error], false},

		// empty interface
		{"any", IsInterface[any], true},
		// built-in interface
		{"error", IsInterface[error], true},
		{"testing.TB", IsInterface[testing.TB], true},

		// custom interfaces
		{"myInterface", IsInterface[myInterface], true},
		{"embedded interface", IsInterface[embedded], true},
		{"interface{ Foo() }", IsInterface[interface{ Foo() }], true},

		// union of interface types
		{
			"multiple methods interface",
			IsInterface[interface {
				String() string
				Error() string
			}],
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// call generic fn by asserting to correct type
			switch fn := tt.fn.(type) {
			case func() bool:
				if got := fn(); got != tt.expected {
					t.Fatalf("expected %v, got %v", tt.expected, got)
				}
			default:
				t.Fatalf("invalid test fn")
			}
		})
	}
}
