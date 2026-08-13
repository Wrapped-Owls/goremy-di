package remychef

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

type fixedHealthchecker struct{ err error }

func (f fixedHealthchecker) HealthCheck() error { return f.err }

// blocks once started, so a test can count how many checks run at the same time
type gatedHealthchecker struct {
	started *chan struct{}
	release *chan struct{}
}

func (g gatedHealthchecker) HealthCheck(ctx context.Context) error {
	*g.started <- struct{}{}
	<-*g.release
	return nil
}

func TestApp_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		values []any
		want   map[string]error
	}{
		{
			name:   "no tracked services",
			values: nil,
			want:   map[string]error{},
		},
		{
			name:   "non-checker services are skipped",
			values: []any{"plain string", 42},
			want:   map[string]error{},
		},
		{
			name: "mixed healthy and failing",
			values: []any{
				fixedHealthchecker{err: nil},
				fixedHealthchecker{err: errBoom},
			},
			want: map[string]error{
				"remychef.fixedHealthchecker":   nil,
				"remychef.fixedHealthchecker#2": errBoom,
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector())
			for _, value := range testCase.values {
				app.track(value)
			}

			got := app.HealthCheck(context.Background())
			if len(got) != len(testCase.want) {
				t.Fatalf("len(got) = %d, want %d (%+v)", len(got), len(testCase.want), got)
			}
			for name, wantErr := range testCase.want {
				if !errors.Is(got[name], wantErr) {
					t.Fatalf("got[%q] = %v, want %v", name, got[name], wantErr)
				}
			}
		})
	}
}

func TestApp_HealthCheck_boundedParallelism(t *testing.T) {
	t.Parallel()

	const (
		concurrency = 2
		total       = 5
	)

	app := New(remy.NewInjector(), Config{HealthCheckConcurrency: concurrency})

	started := make(chan struct{}, total)
	release := make(chan struct{})
	for checker := 0; checker < total; checker++ {
		app.track(gatedHealthchecker{started: &started, release: &release})
	}

	resultCh := make(chan map[string]error, 1)
	go func() { resultCh <- app.HealthCheck(context.Background()) }()

	// Exactly `concurrency` checks can be running at once: receiving this many
	// `started` signals deterministically proves that many (and no more, since a
	// (concurrency+1)-th checker would block acquiring the semaphore, never sending).
	for awaited := 0; awaited < concurrency; awaited++ {
		<-started
	}

	select {
	case <-started:
		t.Fatalf("more than %d checks started before any was released", concurrency)
	default:
	}

	close(release)

	result := <-resultCh
	if len(result) != total {
		t.Fatalf("len(result) = %d, want %d", len(result), total)
	}
}

func TestApp_HealthCheck_timeout(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector(), Config{HealthCheckTimeout: -1})
	var ranWithDeadline sync.Once
	var sawErr error

	app.track(healthcheckerFunc(func(ctx context.Context) error {
		ranWithDeadline.Do(func() { sawErr = ctx.Err() })
		return ctx.Err()
	}))

	got := app.HealthCheck(context.Background())
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
	for _, err := range got {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want %v", err, context.DeadlineExceeded)
		}
	}
	if !errors.Is(sawErr, context.DeadlineExceeded) {
		t.Fatalf("ctx passed to check = %v, want already expired", sawErr)
	}
}

type healthcheckerFunc func(context.Context) error

func (f healthcheckerFunc) HealthCheck(ctx context.Context) error { return f(ctx) }
