package remy

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

var errAuditFailed = errors.New("audit sink refused the language")

func languageBind() Bind[fixtures.Language] {
	return Factory(
		func(DependencyRetriever) (fixtures.Language, error) {
			return fixtures.GoProgrammingLang{}, nil
		},
	)
}

func auditLanguage(
	retriever DependencyRetriever, language fixtures.Language,
) (fixtures.Language, error) {
	audit, err := Get[*fixtures.AuditLog](retriever)
	if err != nil {
		return nil, err
	}

	return fixtures.AuditedLanguage{Inner: language, Audit: audit}, nil
}

func TestDecorate_AddsAuditingWithoutChangingTheWrappedBind(t *testing.T) {
	t.Parallel()

	audit := &fixtures.AuditLog{}
	injector := NewInjector()
	RegisterInstance(injector, audit)
	Register(injector, Decorate[fixtures.Language](languageBind(), auditLanguage))

	language, err := Get[fixtures.Language](injector)
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}

	if name := language.Name(); name != "Go" {
		t.Fatalf("Name() = %q, want the answer of the bind that was wrapped", name)
	}
	if kind := language.Kind(); kind != "programming" {
		t.Fatalf("Kind() = %q, want the untouched answer to pass through", kind)
	}
	if len(audit.Lookups) != 1 || audit.Lookups[0] != "Go" {
		t.Fatalf("audit holds %v, want the decorator to have recorded one Go lookup", audit.Lookups)
	}
}

func TestDecorate_ResolvesTheDecoratorCollaboratorFromTheInjector(t *testing.T) {
	t.Parallel()

	injector := NewInjector()
	Register(injector, Decorate[fixtures.Language](languageBind(), auditLanguage))

	if _, err := Get[fixtures.Language](injector); err == nil {
		t.Fatal("Get() error = nil, want the unresolved audit log to fail the decoration")
	}
}

func TestDecorate_StacksFromTheInsideOut(t *testing.T) {
	t.Parallel()

	injector := NewInjector()
	inner := Decorate[string](Instance("go"), suffixDecorator("-inner"))
	Register(injector, Decorate[string](inner, suffixDecorator("-outer")))

	got, err := Get[string](injector)
	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if got != "go-inner-outer" {
		t.Fatalf("Get() = %q, want each decorator applied around the previous one", got)
	}
}

func TestDecorate_SurfacesTheDecoratorError(t *testing.T) {
	t.Parallel()

	injector := NewInjector()
	Register(
		injector, Decorate[fixtures.Language](
			languageBind(), func(
				DependencyRetriever, fixtures.Language,
			) (fixtures.Language, error) {
				return nil, errAuditFailed
			},
		),
	)

	if _, err := Get[fixtures.Language](injector); !errors.Is(err, errAuditFailed) {
		t.Fatalf("Get() error = %v, want %v", err, errAuditFailed)
	}
}

func TestDecorate_RunsPerBindKind(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		bind            func(builds *int) Bind[string]
		wantBuilds      int
		wantDecorations int
	}{
		{
			name:            "an instance is decorated once, while it registers",
			bind:            func(*int) Bind[string] { return Instance("go") },
			wantBuilds:      0,
			wantDecorations: 1,
		},
		{
			name: "a factory decorates every value it builds",
			bind: func(builds *int) Bind[string] {
				return Factory(countingBinder(builds))
			},
			wantBuilds:      3,
			wantDecorations: 3,
		},
		{
			name: "an eager singleton builds and decorates once, while it registers",
			bind: func(builds *int) Bind[string] {
				return Singleton(countingBinder(builds))
			},
			wantBuilds:      1,
			wantDecorations: 1,
		},
		{
			name: "a lazy singleton builds and decorates once, on the first retrieval",
			bind: func(builds *int) Bind[string] {
				return LazySingleton(countingBinder(builds))
			},
			wantBuilds:      1,
			wantDecorations: 1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var builds, decorations int
			injector := NewInjector()
			Register(
				injector,
				Decorate[string](testCase.bind(&builds), countingDecorator(&decorations)),
			)

			for round := 0; round < 3; round++ {
				if _, err := Get[string](injector); err != nil {
					t.Fatalf("retrieval %d: unexpected error: %v", round, err)
				}
			}

			if builds != testCase.wantBuilds {
				t.Errorf("the binder ran %d times, want %d", builds, testCase.wantBuilds)
			}
			if decorations != testCase.wantDecorations {
				t.Errorf("the decorator ran %d times, want %d", decorations, testCase.wantDecorations)
			}
		})
	}
}

func suffixDecorator(suffix string) Decorator[string] {
	return func(_ DependencyRetriever, value string) (string, error) {
		return value + suffix, nil
	}
}

func countingBinder(builds *int) Binder[string] {
	return func(DependencyRetriever) (string, error) {
		*builds++
		return "go", nil
	}
}

func countingDecorator(decorations *int) Decorator[string] {
	return func(_ DependencyRetriever, value string) (string, error) {
		*decorations++
		return value, nil
	}
}
