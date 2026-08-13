package remychef

import (
	"context"
	"errors"
	"testing"
)

type plainHealthchecker struct{ err error }

func (p plainHealthchecker) HealthCheck() error { return p.err }

type ctxHealthchecker struct{ err error }

func (c ctxHealthchecker) HealthCheck(context.Context) error { return c.err }

type plainShutdowner struct{ err error }

func (p plainShutdowner) Shutdown() error { return p.err }

type ctxShutdowner struct{ err error }

func (c ctxShutdowner) Shutdown(context.Context) error { return c.err }

var errBoom = errors.New("boom")

func TestHealthCheckerOf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		value   any
		wantOK  bool
		wantErr error
	}{
		{
			name:    "no-context checker",
			value:   plainHealthchecker{err: errBoom},
			wantOK:  true,
			wantErr: errBoom,
		},
		{
			name:    "context checker",
			value:   ctxHealthchecker{err: errBoom},
			wantOK:  true,
			wantErr: errBoom,
		},
		{name: "not a checker", value: "plain string", wantOK: false},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			check, ok := healthCheckerOf(testCase.value)
			if ok != testCase.wantOK {
				t.Fatalf("ok = %v, want %v", ok, testCase.wantOK)
			}
			if !ok {
				return
			}
			if err := check(context.Background()); !errors.Is(err, testCase.wantErr) {
				t.Fatalf("err = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}

func TestShutdownerOf(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		value   any
		wantOK  bool
		wantErr error
	}{
		{
			name:    "no-context shutdowner",
			value:   plainShutdowner{err: errBoom},
			wantOK:  true,
			wantErr: errBoom,
		},
		{
			name:    "context shutdowner",
			value:   ctxShutdowner{err: errBoom},
			wantOK:  true,
			wantErr: errBoom,
		},
		{name: "not a shutdowner", value: 42, wantOK: false},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			teardown, ok := shutdownerOf(testCase.value)
			if ok != testCase.wantOK {
				t.Fatalf("ok = %v, want %v", ok, testCase.wantOK)
			}
			if !ok {
				return
			}
			if err := teardown(context.Background()); !errors.Is(err, testCase.wantErr) {
				t.Fatalf("err = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}
