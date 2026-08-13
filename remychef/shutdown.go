package remychef

import (
	"context"
	"errors"
	"fmt"
)

// ShutdownReport aggregates the per-service errors from an App.Shutdown call.
type ShutdownReport struct {
	// Errors maps a service name to the error its Shutdown call returned; services that
	// shut down cleanly, or don't implement Shutdowner/ShutdownerWithContext, are absent.
	Errors map[string]error
}

// Err joins every per-service error into one, or returns nil when Errors is empty.
func (report *ShutdownReport) Err() error {
	if len(report.Errors) == 0 {
		return nil
	}

	joined := make([]error, 0, len(report.Errors))
	for name, err := range report.Errors {
		joined = append(joined, fmt.Errorf("%s: %w", name, err))
	}
	return errors.Join(joined...)
}

// Shutdown tears down every tracked service in reverse construction order,
// aggregating failures instead of stopping at the first one, and returns the same
// report on every later call. remy resolves synchronously and depth-first, so
// construction order is already a valid reverse-topological order.
func (app *App) Shutdown(ctx context.Context) *ShutdownReport {
	app.mu.Lock()
	if app.shutdown != nil {
		report := app.shutdown
		app.mu.Unlock()
		return report
	}
	services := app.built
	app.built = nil
	app.mu.Unlock()

	failures := make(map[string]error)
	for index := len(services) - 1; index >= 0; index-- {
		service := services[index]
		teardown, ok := shutdownerOf(service.value)
		if !ok {
			continue
		}
		if err := teardown(ctx); err != nil {
			failures[service.name] = err
		}
	}

	report := &ShutdownReport{Errors: failures}
	app.mu.Lock()
	app.shutdown = report
	app.mu.Unlock()
	return report
}
