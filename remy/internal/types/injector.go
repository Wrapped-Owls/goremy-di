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
		GetNamed(T, string) (any, bool)
		Get(T) (any, bool)
		CheckReflectionOptions
	}
	Storage[T comparable] interface {
		ValuesSetter[T]
		ValuesGetter[T]
	}
	DependencyRetriever interface {
		RetrieveBind(BindKey) (any, bool)
		RetrieveNamedBind(string, BindKey) (any, bool)
		ValuesGetter[BindKey]
	}
	Injector interface {
		Bind(BindKey, any)
		BindNamed(string, BindKey, any)
		ValuesSetter[BindKey]
		DependencyRetriever
	}
)
