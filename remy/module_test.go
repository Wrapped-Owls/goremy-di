package remy

import (
	"errors"
	"testing"
)

func TestModule_BasicRegistration(t *testing.T) {
	inj := NewInjector()

	mod := NewModule(
		WithInstance("hello"),
		WithBind(LazySingleton(func(_ DependencyRetriever) (int, error) { return 42, nil })),
	)

	if err := RegisterModule(inj, mod); err != nil {
		t.Fatalf("unexpected error registering module: %v", err)
	}

	if got := MustGet[string](inj); got != "hello" {
		t.Fatalf("unexpected string value: %q", got)
	}
	if got := MustGet[int](inj); got != 42 {
		t.Fatalf("unexpected int value: %d", got)
	}
}

func TestModule_Constructors(t *testing.T) {
	inj := NewInjector()

	// Provide a base dependency
	base := NewModule(
		WithInstance("world"),
	)

	// Build value from constructor with no args
	type NoArg string
	consNoArg := NewModule(
		WithConstructor(Factory[NoArg], func() (NoArg, error) { return "no-arg", nil }),
	)

	// Build value from constructor with one arg (pulls string from injector)
	type FromStrLen int
	consOneArg := NewModule(
		WithConstructor1(Factory[FromStrLen], func(s string) (FromStrLen, error) {
			return FromStrLen(len(s)), nil
		}),
	)

	if err := RegisterModule(inj, base, consNoArg, consOneArg); err != nil {
		t.Fatalf("unexpected error registering modules: %v", err)
	}

	if got := MustGet[NoArg](inj); string(got) != "no-arg" {
		t.Fatalf("unexpected constructed NoArg: %q", got)
	}

	if got := MustGet[FromStrLen](inj); int(got) != len("world") {
		t.Fatalf("unexpected constructed FromStrLen: %d", got)
	}
}

func TestRegisterModule_RecoveryAndPartialApply(t *testing.T) {
	inj := NewInjector() // default: cannot override

	mod := NewModule(
		WithInstance(1),           // ok
		WithInstance(2),           // duplicate int, should panic
		WithInstance("unreached"), // should not execute after panic
	)

	err := RegisterModule(inj, mod)
	if err == nil {
		t.Fatalf("expected error from RegisterModule, got nil")
	}

	// Ensure some error was returned
	if !errors.Is(err, ErrAlreadyBound) {
		t.Fatalf("unexpected empty error returned")
	}

	// First registration should remain applied
	if got := MustGet[int](inj); got != 1 {
		t.Fatalf("expected first int registration to persist, got %d", got)
	}

	// The third registration should not have been applied
	if _, e := Get[string](inj); e == nil {
		t.Fatalf("expected string to be absent due to early panic, but it was registered")
	}
}

func TestRegisterModule_ComposeMultiple(t *testing.T) {
	inj := NewInjector()

	m1 := NewModule(WithInstance(uint8(7)))
	m2 := NewModule(
		WithBind(Factory(func(_ DependencyRetriever) (bool, error) { return true, nil })),
	)

	if err := RegisterModule(inj, m1, m2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if MustGet[uint8](inj) != 7 {
		t.Fatalf("unexpected uint8 value")
	}
	if !MustGet[bool](inj) {
		t.Fatalf("unexpected bool value")
	}
}
