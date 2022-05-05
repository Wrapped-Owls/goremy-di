package types

type (
	Storage interface {
		Set(value any)
		SetNamed(key string, value any)
		Get(key string) any
		Binds() []any
	}
	DependencyRetriever interface {
		Get() any
		GetNamed() any
	}
	Injector interface {
		Storage() Storage
		Bind(BindKey, any)
		RetrieveBind(BindKey) (any, bool)
		DependencyRetriever
	}
)
