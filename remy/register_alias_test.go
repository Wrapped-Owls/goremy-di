package remy

import (
	"errors"
	"testing"
)

type aliasGreeter interface{ Greet() string }

type aliasEnglish struct{ suffix string }

func (e *aliasEnglish) Greet() string { return "hello" + e.suffix }

func TestRegisterAs(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(inj Injector)
		tag      []string
		wantErr  error
		wantWord string
	}{
		{
			name: "factory alias resolves concrete implementation",
			setup: func(inj Injector) {
				RegisterInstance(inj, &aliasEnglish{})
				RegisterAs(
					inj,
					Factory[aliasGreeter],
					func(e *aliasEnglish) aliasGreeter { return e },
				)
			},
			wantWord: "hello",
		},
		{
			name: "singleton alias resolves eagerly at registration",
			setup: func(inj Injector) {
				RegisterInstance(inj, &aliasEnglish{suffix: "!"})
				RegisterAs(
					inj,
					Singleton[aliasGreeter],
					func(e *aliasEnglish) aliasGreeter { return e },
				)
			},
			wantWord: "hello!",
		},
		{
			name: "tagged alias resolves the tagged concrete bind",
			setup: func(inj Injector) {
				RegisterInstance(inj, &aliasEnglish{suffix: " there"}, "polite")
				RegisterAs(
					inj,
					Factory[aliasGreeter],
					func(e *aliasEnglish) aliasGreeter { return e },
					"polite",
				)
			},
			tag:      []string{"polite"},
			wantWord: "hello there",
		},
		{
			name: "alias without registered concrete fails with not registered",
			setup: func(inj Injector) {
				RegisterAs(
					inj,
					Factory[aliasGreeter],
					func(e *aliasEnglish) aliasGreeter { return e },
				)
			},
			wantErr: ErrElementNotRegistered,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			inj := NewInjector()
			testCase.setup(inj)

			greeter, err := Get[aliasGreeter](inj, testCase.tag...)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("error = %v, want %v", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if word := greeter.Greet(); word != testCase.wantWord {
				t.Fatalf("Greet() = %q, want %q", word, testCase.wantWord)
			}
		})
	}
}
