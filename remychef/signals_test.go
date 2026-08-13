package remychef

import (
	"context"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

// the OS-facing half is a thin signal.NotifyContext wrapper the stdlib already covers
func TestApp_shutdownOnDone(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	var order []string
	app.track(recordingShutdowner{name: "svc", order: &order})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	report := app.shutdownOnDone(ctx)
	if err := report.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	if len(order) != 1 {
		t.Fatalf("len(order) = %d, want 1", len(order))
	}
}
