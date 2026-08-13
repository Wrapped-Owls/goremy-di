package remychef

import (
	"context"
	"os"
	"os/signal"
)

// ShutdownOnSignals blocks until one of sig arrives (os.Interrupt if none given), then
// runs Shutdown with a fresh context.Background() and returns its report.
func (app *App) ShutdownOnSignals(sig ...os.Signal) *ShutdownReport {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt}
	}

	ctx, stop := signal.NotifyContext(context.Background(), sig...)
	defer stop()

	return app.shutdownOnDone(ctx)
}

// split out so tests drive it with a cancelable context instead of a real OS signal
func (app *App) shutdownOnDone(ctx context.Context) *ShutdownReport {
	<-ctx.Done()
	return app.Shutdown(context.Background())
}
