package remychef

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
)

func TestWithConstructor_arities(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		register func(app *App) ModuleRegister
		want     string
	}{
		{
			name: "one dependency",
			register: func(*App) ModuleRegister {
				return WithConstructor1(remy.LazySingleton[widget],
					func(first string) (widget, error) { return widget{name: first}, nil })
			},
			want: "a",
		},
		{
			name: "two dependencies",
			register: func(*App) ModuleRegister {
				return WithConstructor2(remy.LazySingleton[widget],
					func(first string, second int) (widget, error) {
						return widget{name: first + string(rune('0'+second))}, nil
					})
			},
			want: "a2",
		},
		{
			name: "three dependencies",
			register: func(*App) ModuleRegister {
				return WithConstructor3(remy.LazySingleton[widget],
					func(first string, second int, third bool) (widget, error) {
						if !third {
							return widget{}, errors.New("third must be true")
						}
						return widget{name: first + string(rune('0'+second)) + "c"}, nil
					})
			},
			want: "a2c",
		},
		{
			name: "four dependencies",
			register: func(*App) ModuleRegister {
				return WithConstructor4(remy.LazySingleton[widget],
					func(first string, second int, third bool, fourth uint8) (widget, error) {
						return widget{
							name: first + string(rune('0'+second)) + "c" + string(rune('0'+fourth)),
						}, nil
					})
			},
			want: "a2c3",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector())
			remy.RegisterInstance(app.Injector, "a")
			remy.RegisterInstance(app.Injector, 2)
			remy.RegisterInstance(app.Injector, true)
			remy.RegisterInstance(app.Injector, uint8(3))

			if err := app.Register(testCase.register(app)); err != nil {
				t.Fatalf("Register: %v", err)
			}

			got := remy.MustGet[widget](app.Injector)
			if got.name != testCase.want {
				t.Fatalf("name = %q, want %q", got.name, testCase.want)
			}
			if len(app.built) != 1 {
				t.Fatalf("len(built) = %d, want the built widget tracked once", len(app.built))
			}
		})
	}
}

func TestWithConstructor_tracksPerBindKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		bindFunc  func(remy.Binder[widget]) remy.Bind[widget]
		retrieve  int
		wantBuilt int
	}{
		{
			name:      "lazy singleton is tracked once",
			bindFunc:  remy.LazySingleton[widget],
			retrieve:  3,
			wantBuilt: 1,
		},
		{
			name:      "singleton is tracked at registration",
			bindFunc:  remy.Singleton[widget],
			retrieve:  0,
			wantBuilt: 1,
		},
		{
			name:      "factory is never tracked",
			bindFunc:  remy.Factory[widget],
			retrieve:  3,
			wantBuilt: 0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			app := New(remy.NewInjector())
			remy.RegisterInstance(app.Injector, "a")
			register := WithConstructor1(testCase.bindFunc,
				func(first string) (widget, error) { return widget{name: first}, nil })
			if err := app.Register(register); err != nil {
				t.Fatalf("Register: %v", err)
			}

			for index := 0; index < testCase.retrieve; index++ {
				remy.MustGet[widget](app.Injector)
			}

			if len(app.built) != testCase.wantBuilt {
				t.Fatalf("len(built) = %d, want %d", len(app.built), testCase.wantBuilt)
			}
		})
	}
}

func TestWithConstructor_error(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("boom")
	app := New(remy.NewInjector())
	remy.RegisterInstance(app.Injector, "a")
	register := WithConstructor1(remy.LazySingleton[widget],
		func(string) (widget, error) { return widget{}, wantErr })
	if err := app.Register(register); err != nil {
		t.Fatalf("Register: %v", err)
	}

	if _, err := remy.Get[widget](app.Injector); !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
	if len(app.built) != 0 {
		t.Fatalf("len(built) = %d, want a failed build to stay untracked", len(app.built))
	}
}
