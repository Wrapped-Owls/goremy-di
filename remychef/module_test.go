package remychef

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestApp_RegisterModule(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		modules []remy.Module
		wantErr bool
	}{
		{
			name:    "single module registers",
			modules: []remy.Module{remy.NewModule(remy.WithInstance("hello"))},
			wantErr: false,
		},
		{
			name: "duplicate registration errors",
			modules: []remy.Module{
				remy.NewModule(remy.WithInstance(1)),
				remy.NewModule(remy.WithInstance(2)),
			},
			wantErr: true,
		},
	}
	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector())
			err := app.RegisterModule(testCase.modules...)
			if testCase.wantErr && !errors.Is(err, remy.ErrAlreadyBound) {
				t.Fatalf("err = %v, want wrapping %v", err, remy.ErrAlreadyBound)
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("err = %v, want nil", err)
			}
		})
	}
}

func TestApp_Module(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	module := app.Module(
		WithInstance("from-module"),
		WithLazySingleton(func(remy.DependencyRetriever) (widget, error) {
			return widget{name: "lazy"}, nil
		}),
		WithFactory(func(remy.DependencyRetriever) (int, error) { return 7, nil }),
	)

	if err := app.RegisterModule(module); err != nil {
		t.Fatalf("RegisterModule: %v", err)
	}

	// the instance is tracked at registration, the lazy one only once retrieved
	if len(app.built) != 1 {
		t.Fatalf("len(built) = %d, want only the instance tracked", len(app.built))
	}
	if got := remy.MustGet[widget](app.Injector); got.name != "lazy" {
		t.Fatalf("widget = %q, want %q", got.name, "lazy")
	}
	for retrieval := 0; retrieval < 3; retrieval++ {
		remy.MustGet[int](app.Injector)
	}
	if len(app.built) != 2 {
		t.Fatalf("len(built) = %d, want the factory to stay untracked", len(app.built))
	}
	if got := remy.MustGet[string](app.Injector); got != "from-module" {
		t.Fatalf("string = %q, want %q", got, "from-module")
	}
}

func TestApp_RegisterModule_bypassesTracking(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	mod := remy.NewModule(remy.WithSingleton(func(remy.DependencyRetriever) (widget, error) {
		return widget{name: "raw"}, nil
	}))

	if err := app.RegisterModule(mod); err != nil {
		t.Fatalf("RegisterModule: %v", err)
	}
	if len(app.built) != 0 {
		t.Fatalf(
			"len(built) = %d, want 0 (raw remy.Module registrations bypass App tracking)",
			len(app.built),
		)
	}
}
