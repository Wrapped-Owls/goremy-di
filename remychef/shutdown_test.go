package remychef

import (
	"context"
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

type recordingShutdowner struct {
	name  string
	order *[]string
	err   error
}

func (r recordingShutdowner) Shutdown() error {
	*r.order = append(*r.order, r.name)
	return r.err
}

func TestShutdownReport_Err(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		errs    map[string]error
		wantNil bool
	}{
		{name: "no errors", errs: map[string]error{}, wantNil: true},
		{name: "one error", errs: map[string]error{"svc": errBoom}, wantNil: false},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			report := &ShutdownReport{Errors: testCase.errs}
			err := report.Err()
			if testCase.wantNil {
				if err != nil {
					t.Fatalf("Err() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatal("Err() = nil, want an error")
			}
			if !errors.Is(err, errBoom) {
				t.Fatalf("Err() = %v, want it to wrap %v", err, errBoom)
			}
		})
	}
}

func TestApp_Shutdown_reverseOrder(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	var order []string
	app.track(recordingShutdowner{name: "first", order: &order})
	app.track(recordingShutdowner{name: "second", order: &order})
	app.track("not a shutdowner")
	app.track(recordingShutdowner{name: "third", order: &order})

	report := app.Shutdown(context.Background())
	if err := report.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}

	want := []string{"third", "second", "first"}
	if len(order) != len(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
	for index, name := range want {
		if order[index] != name {
			t.Fatalf("order[%d] = %q, want %q", index, order[index], name)
		}
	}
}

func TestApp_Shutdown_aggregatesErrors(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	var order []string
	app.track(recordingShutdowner{name: "a", order: &order, err: errBoom})
	app.track(recordingShutdowner{name: "b", order: &order, err: errBoom})

	report := app.Shutdown(context.Background())
	if len(report.Errors) != 2 {
		t.Fatalf(
			"len(Errors) = %d, want 2 (both must run despite the first failing)",
			len(report.Errors),
		)
	}
	if len(order) != 2 {
		t.Fatalf("len(order) = %d, want 2", len(order))
	}
}

func TestApp_Shutdown_idempotent(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	var order []string
	app.track(recordingShutdowner{name: "once", order: &order})

	first := app.Shutdown(context.Background())
	second := app.Shutdown(context.Background())

	if first != second {
		t.Fatalf("second Shutdown() returned a different report")
	}
	if len(order) != 1 {
		t.Fatalf("len(order) = %d, want 1 (Shutdown must not re-run)", len(order))
	}
}
