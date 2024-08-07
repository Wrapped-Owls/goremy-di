package injector

import (
	"errors"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func Register[T any](ij types.Injector, bind types.Bind[T], keys ...string) error {
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	elementType := utils.GetKey[T](injopts.KeyOptsFromStruct(ij.ReflectOpts()))
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

func checkSavedAsBind[T any](
	retriever types.DependencyRetriever, checkElem any,
) (foundElem *T, err error) {
	if genericBind, assertOk := checkElem.(interface {
		PointerValue() any
		GenAsAny(injector types.DependencyRetriever) (any, error)
	}); assertOk {
		// Check if the returned value can implement the requested interface
		if _, ok := genericBind.PointerValue().(T); !ok {
			return
		}
		var anyVal any
		if anyVal, err = genericBind.GenAsAny(retriever); err != nil {
			return
		} else if bindElem, ok := anyVal.(T); ok {
			foundElem = &bindElem
		}
	}
	return
}

func GetAll[T any](
	retriever types.DependencyRetriever, optKey ...string,
) (resultList []T, err error) {
	var elementList []any
	if elementList, err = retriever.GetAll(optKey...); err != nil {
		return
	}

	resultList = make([]T, 0, len(elementList))
	for _, checkElem := range elementList {
		switch instanceBind := checkElem.(type) {
		case T:
			resultList = append(resultList, instanceBind)
		default:
			var foundElem *T
			if foundElem, err = checkSavedAsBind[T](retriever, checkElem); err != nil {
				return
			}

			if foundElem != nil {
				resultList = append(resultList, *foundElem)
			}
		}
	}

	if len(resultList) == 0 {
		err = utils.ErrElementNotRegistered
	}

	return
}

func getByGuess[T any](
	retriever types.DependencyRetriever, optKey ...string,
) (element T, err error) {
	var elementList []T
	if elementList, err = GetAll[T](retriever, optKey...); err != nil {
		return
	}

	totalFound := len(elementList)
	if totalFound == 1 {
		element = elementList[0]
		return
	}

	err = utils.ErrFoundMoreThanOneValidDI
	if totalFound == 0 {
		err = utils.ErrElementNotRegistered
	}

	return
}

func Get[T any](retriever types.DependencyRetriever, keys ...string) (element T, err error) {
	var (
		key         string
		bind        any
		elementType = utils.GetKey[T](injopts.KeyOptsFromStruct(retriever.ReflectOpts()))
	)

	if len(keys) > 0 {
		key = keys[0]
	}
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

	// Start to search for every element if it is configured in this way
	foundElement, accessAllError := getByGuess[T](retriever, keys...)
	if accessAllError == nil {
		element = foundElement
		err = nil
	} else if !errors.Is(accessAllError, utils.ErrElementNotRegistered) {
		err = accessAllError
	}

	// retrieve values from cacheStorage
	return
}

func TryGet[T any](retriever types.DependencyRetriever, keys ...string) (result T) {
	result, _ = Get[T](retriever, keys...)
	return
}

func GetGen[T any](
	retriever types.DependencyRetriever, elements []types.InstancePair[any], keys ...string,
) (result T, err error) {
	subInjector := New(injopts.CacheOptNone, retriever.ReflectOpts(), retriever)
	for _, element := range elements {
		var (
			opts       = injopts.KeyOptsFromStruct(subInjector.ReflectOpts())
			typeSeeker = element.Value
		)
		if element.InterfaceValue != nil {
			opts |= injopts.KeyOptIgnorePointer
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
	subInjector := New(injopts.CacheOptNone, retriever.ReflectOpts(), retriever)
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
