package types

type (
	valuesGetter interface {
		Get(key string) any
		Binds() []any
	}
	Storage interface {
		Set(value any)
		SetNamed(key string, value any)
		valuesGetter
	}
	DependencyRetriever interface {
		RetrieveBind(BindKey) (any, bool)
		valuesGetter
	}
	Injector interface {
		Storage() Storage
		Bind(BindKey, any)
		DependencyRetriever
	}
)
