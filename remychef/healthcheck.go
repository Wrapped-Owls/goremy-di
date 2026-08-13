package remychef

import (
	"context"
	"sync"
)

// HealthCheck probes every tracked service implementing Healthchecker or
// HealthcheckerWithContext, at most Config.HealthCheckConcurrency at a time, each
// capped by Config.HealthCheckTimeout. It reports one entry per probed service.
func (app *App) HealthCheck(ctx context.Context) map[string]error {
	app.mu.Lock()
	services := make([]trackedService, len(app.built))
	copy(services, app.built)
	app.mu.Unlock()

	results := make(map[string]error)
	var (
		resultsMu sync.Mutex
		running   sync.WaitGroup
	)
	slots := make(chan struct{}, app.healthCheckConcurrency())

	for _, service := range services {
		check, ok := healthCheckerOf(service.value)
		if !ok {
			continue
		}

		running.Add(1)
		slots <- struct{}{}
		// go.mod pins 1.20: the loop variables must be passed in, not captured
		go func(name string, check func(context.Context) error) {
			defer running.Done()
			defer func() { <-slots }()

			checkCtx, cancel := context.WithTimeout(ctx, app.healthCheckTimeout())
			defer cancel()

			err := check(checkCtx)

			resultsMu.Lock()
			results[name] = err
			resultsMu.Unlock()
		}(service.name, check)
	}

	running.Wait()
	return results
}
