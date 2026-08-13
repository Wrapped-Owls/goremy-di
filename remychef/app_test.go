package remychef

import (
	"testing"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		cfg  []Config
		want Config
	}{
		{name: "zero value uses defaults", cfg: nil, want: Config{}},
		{
			name: "explicit config kept as-is",
			cfg:  []Config{{HealthCheckTimeout: time.Second, HealthCheckConcurrency: 3}},
			want: Config{HealthCheckTimeout: time.Second, HealthCheckConcurrency: 3},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			injector := remy.NewInjector()
			app := New(injector, testCase.cfg...)

			if app.Injector != injector {
				t.Fatalf("Injector = %v, want %v", app.Injector, injector)
			}
			if app.cfg != testCase.want {
				t.Fatalf("cfg = %+v, want %+v", app.cfg, testCase.want)
			}
		})
	}
}

func TestApp_track(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	app.track("first")
	app.track("second")
	app.track(7)

	if len(app.built) != 3 {
		t.Fatalf("len(built) = %d, want 3", len(app.built))
	}

	wantNames := []string{"string", "string#2", "int"}
	for index, want := range wantNames {
		if got := app.built[index].name; got != want {
			t.Fatalf("built[%d].name = %q, want %q", index, got, want)
		}
	}
}

func TestApp_healthCheckTimeout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		cfg  Config
		want time.Duration
	}{
		{name: "zero uses default", cfg: Config{}, want: DefaultHealthCheckTimeout},
		{name: "positive kept", cfg: Config{HealthCheckTimeout: time.Minute}, want: time.Minute},
		{name: "negative kept for tests", cfg: Config{HealthCheckTimeout: -1}, want: -1},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector(), testCase.cfg)
			if got := app.healthCheckTimeout(); got != testCase.want {
				t.Fatalf("healthCheckTimeout() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestApp_healthCheckConcurrency(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		cfg  Config
		want int
	}{
		{name: "zero uses default", cfg: Config{}, want: DefaultHealthCheckConcurrency},
		{name: "positive kept", cfg: Config{HealthCheckConcurrency: 4}, want: 4},
		{name: "negative floored at one", cfg: Config{HealthCheckConcurrency: -5}, want: 1},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector(), testCase.cfg)
			if got := app.healthCheckConcurrency(); got != testCase.want {
				t.Fatalf("healthCheckConcurrency() = %d, want %d", got, testCase.want)
			}
		})
	}
}
