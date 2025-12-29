package types

type (
	BindKey interface {
		comparable()
	}

	BindDependencies[T any] map[BindKey]T
	DependencyGraph         struct {
		UnnamedDependency BindDependencies[bool]
		NamedDependency   BindDependencies[map[string]bool]
	}

	KeyElem[T any] struct{}
)

func (k KeyElem[T]) comparable() {
	// Just a stub function
}
