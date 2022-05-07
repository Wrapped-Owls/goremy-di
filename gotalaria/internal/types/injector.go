package types

type (
	ValuesGetter interface {
		Get(key string) any
		Binds() []any
	}
	Storage interface {
		AllowOverride(bool)
		Set(value any)
		SetNamed(key string, value any)
		ValuesGetter
	}
	DependencyRetriever interface {
		RetrieveBind(BindKey) (any, bool)
		RetrieveNamedBind(name string, bType BindKey) (any, bool)
		ValuesGetter
	}
	Injector interface {
		Storage() Storage
		Bind(BindKey, any)
		BindNamed(name string, bType BindKey, value any)
		DependencyRetriever
	}
)
