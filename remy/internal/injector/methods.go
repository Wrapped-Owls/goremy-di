package injector

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/keyopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func Register[T any](ij types.Injector, bind types.Bind[T], keys ...string) error {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	elementType := utils.GetKey[T](keyopts.FromReflectOpts(ij.ReflectOpts()))
	var retriever types.DependencyRetriever = ij
	if wrappedRetriever := retriever.WrapRetriever(); wrappedRetriever != nil {
		retriever = wrappedRetriever
	}

	var (
		value any = bind
		err   error
	)
	if bindType := bind.Type(); bindType == types.BindInstance || bindType == types.BindSingleton {
		if value, err = bind.Generates(retriever); err != nil {
			return err
		}
	}

	if key != "" {
		return ij.BindNamed(elementType, key, value)
	}
	return ij.Bind(elementType, value)
}

func Get[T any](retriever types.DependencyRetriever, keys ...string) (T, error) {
	var key string

	if len(keys) > 0 {
		key = keys[0]
	}
	elementType := utils.GetKey[T](keyopts.FromReflectOpts(retriever.ReflectOpts()))

	var (
		bind any
		err  error
	)

	if wrappedRetriever := retriever.WrapRetriever(); wrappedRetriever != nil {
		retriever = wrappedRetriever
	}
	// search in dynamic injections that needed to run a given function
	if key != "" {
		bind, err = retriever.GetNamed(elementType, key)
	} else {
		bind, err = retriever.Get(elementType)
	}

	if err == nil {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			return typedBind.Generates(retriever)
		}
		if instanceBind, assertOk := bind.(T); assertOk {
			return instanceBind, nil
		}
		err = utils.ErrTypeCastInRuntime
	}
	// retrieve values from cacheStorage
	return utils.Default[T](), err
}

func TryGet[T any](retriever types.DependencyRetriever, keys ...string) (result T) {
	result, _ = Get[T](retriever, keys...)
	return
}

func GetGen[T any](
	retriever types.DependencyRetriever, elements []types.InstancePair[any], keys ...string,
) (result T, err error) {
	subInjector := New(false, retriever.ReflectOpts(), retriever)
	for _, element := range elements {
		var (
			opts       = keyopts.FromReflectOpts(subInjector.ReflectOpts())
			typeSeeker = element.Value
		)
		if element.InterfaceValue != nil {
			opts |= keyopts.KeyOptIgnorePointer
			typeSeeker = element.InterfaceValue
		}
		bindKey := utils.GetElemKey(typeSeeker, opts)

		if element.Key != "" {
			if err = subInjector.BindNamed(bindKey, element.Key, element.Value); err != nil {
				return
			}
		} else if err = subInjector.Bind(bindKey, element.Value); err != nil {
			return
		}
	}

	return Get[T](subInjector, keys...)
}

func TryGetGen[T any](
	retriever types.DependencyRetriever, elements []types.InstancePair[any], keys ...string,
) (result T) {
	result, _ = GetGen[T](retriever, elements, keys...)
	return
}

func GetGenFunc[T any](
	retriever types.DependencyRetriever,
	binder func(injector types.Injector) error, keys ...string,
) (result T, err error) {
	subInjector := New(false, retriever.ReflectOpts(), retriever)
	if err = binder(subInjector); err != nil {
		return
	}
	return Get[T](subInjector, keys...)
}

func TryGetGenFunc[T any](
	retriever types.DependencyRetriever,
	binder func(injector types.Injector) error,
	keys ...string,
) (result T) {
	result, _ = GetGenFunc[T](retriever, binder, keys...)
	return
}
