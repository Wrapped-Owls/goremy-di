package types

type (
	Storage interface {
		Set(value any)
		SetNamed(key string, value any)
		Get(key string) any
		Binds() []any
	}
	Injector interface {
		Storage() Storage
	}
	DependencyRetriever interface {
		Get() any
		GetNamed() any
	}
)
