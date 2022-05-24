package utils

import "testing"

func TestTypeName__Generify(t *testing.T) {
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

	if TypeName[super](false) == TypeName[sub](false) {
		t.Error("type names was the same when should not generify")
	}

	if TypeName[super](true) != TypeName[sub](true) {
		t.Error("generified type name should be the same")
	}

}
