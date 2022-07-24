package types

type (
	CheckReflectionOptions interface {
		ReflectOpts() ReflectionOptions
	}
	ValuesSetter[T comparable] interface {
		Set(T, any)
		SetNamed(T, string, any)
		CheckReflectionOptions
	}
	ValuesGetter[T comparable] interface {
		GetNamed(T, string) (any, error)
		Get(T) (any, error)
		CheckReflectionOptions
	}
	Storage[T comparable] interface {
		ValuesSetter[T]
		ValuesGetter[T]
	}
	DependencyRetriever interface {
		RetrieveBind(BindKey) (any, error)
		RetrieveNamedBind(string, BindKey) (any, error)
		ValuesGetter[BindKey]
	}
	Injector interface {
		Bind(BindKey, any)
		BindNamed(string, BindKey, any)
		ValuesSetter[BindKey]
		DependencyRetriever
	}
)
