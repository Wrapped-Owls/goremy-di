package types

type (
	ValuesGetter interface {
		Get(key string) any
		Binds() []any
	}
	Storage interface {
		Set(value any)
		SetNamed(key string, value any)
		ValuesGetter
	}
	DependencyRetriever interface {
		RetrieveBind(BindKey) (any, bool)
		ValuesGetter
	}
	Injector interface {
		Storage() Storage
		Bind(BindKey, any)
		DependencyRetriever
	}
)
