package utils

import (
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	aTypes "github.com/wrapped-owls/goremy-di/remy/test/fixtures/a/testtypes"
	bTypes "github.com/wrapped-owls/goremy-di/remy/test/fixtures/b/testtypes"
)

func TestGetKey__Generify(t *testing.T) {
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

	options := injopts.KeyOptNone
	if GetKey[super](options) == GetKey[sub](options) {
		t.Error("type names was the same when should not generify")
	}

	options = injopts.KeyOptGenerifyInterface
	if GetKey[super](options) != GetKey[sub](options) {
		t.Error("generified type name should be the same")
	}
}

func TestGetKey__SameStructWithDifferentPackage(t *testing.T) {
	options := injopts.KeyOptUseReflectionType
	if GetKey[aTypes.Syringe](options) == GetKey[bTypes.Syringe](options) {
		t.Error(
			"type names was the same, when it should be different, because of different packages",
		)
	}

	options = injopts.KeyOptUseReflectionType
	elemKey, err := GetElemKey(t, options)
	if err != nil {
		t.Fatal(err)
	}
	if elemKey != GetKey[*testing.T](options) {
		t.Error("element type should be the same from type and object")
	}
}

func TestGetKey__Functions(t *testing.T) {
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

	optionsCases := [...]injopts.KeyGenOption{
		injopts.KeyOptNone, injopts.KeyOptUseReflectionType,
	}

	for _, optCase := range optionsCases {
		// Check for named and unnamed functions
		t.Run(
			fmt.Sprintf("NamedUnnamed functions - %+v", optCase), func(t *testing.T) {
				cases := [...][2]types.BindKey{
					{GetKey[namedStringCallback](optCase), GetKey[stringCallback](optCase)},
					{GetKey[namedMultiArgsCallback](optCase), GetKey[multiArgsCallback](optCase)},
					{
						GetKey[namedBoolCheckerCallback](optCase),
						GetKey[boolCheckerCallback](optCase),
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
			fmt.Sprintf("Different function types - %+v", optCase), func(t *testing.T) {
				cases := [...][2]types.BindKey{
					{GetKey[namedStringCallback](optCase), GetKey[voidCallback](optCase)},
					{GetKey[stringCallback](optCase), GetKey[multiArgsCallback](optCase)},
					{GetKey[voidCallback](optCase), GetKey[boolCheckerCallback](optCase)},
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
		if GetKey[namedStringCallback](optCase) == GetKey[*namedStringCallback](optCase) {
			t.Error("Function pointer should be different than function type")
		}
	}
}
