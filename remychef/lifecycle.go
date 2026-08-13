package remychef

import "context"

type (
	// Healthchecker reports whether a built service is healthy.
	Healthchecker interface {
		HealthCheck() error
	}

	// HealthcheckerWithContext reports whether a built service is healthy, honoring ctx.
	HealthcheckerWithContext interface {
		HealthCheck(ctx context.Context) error
	}

	// Shutdowner releases a built service's resources.
	Shutdowner interface {
		Shutdown() error
	}

	// ShutdownerWithContext releases a built service's resources, honoring ctx.
	ShutdownerWithContext interface {
		Shutdown(ctx context.Context) error
	}
)

// healthCheckerOf type-switches value into a context-aware health check call, if it has one.
func healthCheckerOf(value any) (func(context.Context) error, bool) {
	switch checker := value.(type) {
	case HealthcheckerWithContext:
		return checker.HealthCheck, true
	case Healthchecker:
		return func(context.Context) error { return checker.HealthCheck() }, true
	default:
		return nil, false
	}
}

// shutdownerOf type-switches value into a context-aware shutdown call, if it has one.
func shutdownerOf(value any) (func(context.Context) error, bool) {
	switch shutdowner := value.(type) {
	case ShutdownerWithContext:
		return shutdowner.Shutdown, true
	case Shutdowner:
		return func(context.Context) error { return shutdowner.Shutdown() }, true
	default:
		return nil, false
	}
}
