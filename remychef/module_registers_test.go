package remychef

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

type widget struct{ name string }

func TestWithLazySingleton(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	if err := app.Register(WithLazySingleton(func(remy.DependencyRetriever) (widget, error) {
		return widget{name: "lazy"}, nil
	})); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if len(app.built) != 0 {
		t.Fatalf(
			"len(built) before Get = %d, want 0 (LazySingleton defers construction)",
			len(app.built),
		)
	}

	got := remy.MustGet[widget](app.Injector)
	if got.name != "lazy" {
		t.Fatalf("got.name = %q, want %q", got.name, "lazy")
	}
	if len(app.built) != 1 {
		t.Fatalf("len(built) after Get = %d, want 1", len(app.built))
	}
}

func TestWithLazySingleton_error(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	wantErr := errors.New("build failed")
	if err := app.Register(WithLazySingleton(func(remy.DependencyRetriever) (widget, error) {
		return widget{}, wantErr
	})); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if _, err := remy.Get[widget](app.Injector); !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
	if len(app.built) != 0 {
		t.Fatalf("len(built) = %d, want 0 (failed build must not track)", len(app.built))
	}
}

func TestWithSingleton(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	if err := app.Register(WithSingleton(func(remy.DependencyRetriever) (widget, error) {
		return widget{name: "eager"}, nil
	})); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if len(app.built) != 1 {
		t.Fatalf("len(built) = %d, want 1 (Singleton builds during Register)", len(app.built))
	}
	if got := app.built[0].value.(widget).name; got != "eager" {
		t.Fatalf("tracked name = %q, want %q", got, "eager")
	}
}

func TestWithFactory(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	if err := app.Register(WithFactory(func(remy.DependencyRetriever) (widget, error) {
		return widget{name: "factory"}, nil
	})); err != nil {
		t.Fatalf("Register: %v", err)
	}

	remy.MustGet[widget](app.Injector)
	remy.MustGet[widget](app.Injector)

	if len(app.built) != 0 {
		t.Fatalf("len(built) = %d, want 0 (Factory instances are never tracked)", len(app.built))
	}
}

func TestWithInstance(t *testing.T) {
	t.Parallel()

	app := New(remy.NewInjector())
	if err := app.Register(WithInstance(widget{name: "instance"})); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if len(app.built) != 1 {
		t.Fatalf("len(built) = %d, want 1", len(app.built))
	}
	if got := remy.MustGet[widget](app.Injector); got.name != "instance" {
		t.Fatalf("got.name = %q, want %q", got.name, "instance")
	}
}
