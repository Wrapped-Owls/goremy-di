package remychef

import (
	"fmt"
	"sync"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"
)

const (
	// DefaultHealthCheckTimeout bounds a single HealthCheck call when Config.HealthCheckTimeout is zero.
	DefaultHealthCheckTimeout = 5 * time.Second

	// DefaultHealthCheckConcurrency bounds parallel checks when Config.HealthCheckConcurrency is zero.
	DefaultHealthCheckConcurrency = 8
)

// Config tunes the App built by New; the zero value is valid and uses the package defaults.
type Config struct {
	// HealthCheckTimeout bounds each individual health check. Zero uses DefaultHealthCheckTimeout;
	// a negative value forces an immediately-expired context, useful in tests.
	HealthCheckTimeout time.Duration

	// HealthCheckConcurrency bounds how many health checks run in parallel.
	// Zero uses DefaultHealthCheckConcurrency; a negative value is treated as 1.
	HealthCheckConcurrency int
}

// kept in construction order, because Shutdown tears services down in reverse
type trackedService struct {
	name  string
	value any
}

// App tracks every instance built through its ModuleRegister helpers so it can
// health-check and shut them down later, without remy needing a hook API of its own.
type App struct {
	// Injector is the wrapped remy.Injector; registering on it directly bypasses
	// tracking, which is what raw remy.Register* calls are for.
	Injector remy.Injector

	cfg Config

	mu       sync.Mutex
	built    []trackedService
	seen     map[string]int
	shutdown *ShutdownReport
}

// New builds an App wrapping the given remy.Injector; cfg is optional, first value wins.
func New(injector remy.Injector, cfg ...Config) *App {
	app := &App{Injector: injector, seen: make(map[string]int)}
	if len(cfg) > 0 {
		app.cfg = cfg[0]
	}
	return app
}

// factory registers deliberately never call this
func (app *App) track(value any) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.built = append(app.built, trackedService{name: app.nameForLocked(value), value: value})
}

// caller must hold app.mu
func (app *App) nameForLocked(value any) string {
	base := fmt.Sprintf("%T", value)
	app.seen[base]++
	if count := app.seen[base]; count > 1 {
		return fmt.Sprintf("%s#%d", base, count)
	}
	return base
}

func (app *App) healthCheckTimeout() time.Duration {
	if app.cfg.HealthCheckTimeout == 0 {
		return DefaultHealthCheckTimeout
	}
	return app.cfg.HealthCheckTimeout
}

func (app *App) healthCheckConcurrency() int {
	switch {
	case app.cfg.HealthCheckConcurrency > 0:
		return app.cfg.HealthCheckConcurrency
	case app.cfg.HealthCheckConcurrency < 0:
		return 1
	default:
		return DefaultHealthCheckConcurrency
	}
}

// wrapping the binder is the whole hook mechanism, so remy needs no hook API
func tracking[T any](app *App, binder remy.Binder[T]) remy.Binder[T] {
	return func(retriever remy.DependencyRetriever) (T, error) {
		value, err := binder(retriever)
		if err != nil {
			return value, err
		}
		app.track(value)
		return value, nil
	}
}

// a free function, not a method, because Go forbids generic methods on *App
func trackResult[T any](app *App) func(T, error) (T, error) {
	return func(value T, err error) (T, error) {
		if err == nil {
			app.track(value)
		}
		return value, err
	}
}

func keepResult[T any](value T, err error) (T, error) { return value, err }
