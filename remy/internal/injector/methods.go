package injector

import (
	"errors"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func Register[T any](ij types.Injector, bind types.Bind[T], optTag ...string) error {
	var tag string
	if len(optTag) > 0 {
		tag = optTag[0]
	}

	bindOpts := types.BindOptions{Tag: tag}
	return registerNewDep[T](ij, bind, bindOpts)
}

func RegisterWithOverride[T any](ij types.Injector, bind types.Bind[T], optTag ...string) error {
	var tag string
	if len(optTag) > 0 {
		tag = optTag[0]
	}

	return registerNewDep[T](ij, bind, types.BindOptions{Tag: tag, SoftOverride: true})
}

func registerNewDep[T any](ij types.Injector, bind types.Bind[T], opts types.BindOptions) error {
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

	return ij.BindElem(elementType, value, opts)
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
		bindKey := utils.GetKey[T](injopts.KeyOptNone)
		err = remyErrs.ErrElementNotRegistered{Key: bindKey}
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

	bindKey := utils.GetKey[T](injopts.KeyOptNone)
	err = remyErrs.ErrMultipleDIDuckTypingCandidates{Type: bindKey, Count: totalFound}
	if totalFound == 0 {
		err = remyErrs.ErrElementNotRegistered{Key: bindKey}
	}

	return
}

func Get[T any](retriever types.DependencyRetriever, tags ...string) (element T, err error) {
	var (
		key         string
		elementType = utils.GetKey[T](injopts.KeyOptsFromStruct(retriever.ReflectOpts()))
	)

	if len(tags) > 0 {
		key = tags[0]
	}
	if wrappedRetriever := retriever.WrapRetriever(); wrappedRetriever != nil {
		retriever = wrappedRetriever
	}

	var bind any
	// search in dynamic injections that needed to run a given function
	if bind, err = retriever.RetrieveBind(elementType, key); err == nil {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			return typedBind.Generates(retriever)
		}
		if instanceBind, assertOk := bind.(T); assertOk {
			return instanceBind, nil
		}
		err = remyErrs.ErrTypeCastInRuntime{ActualValue: bind, Expected: new(T)}
	}

	// Start to search for every element if it is configured in this way
	foundElement, accessAllError := getByGuess[T](retriever, tags...)
	if accessAllError == nil {
		element = foundElement
		err = nil
	} else if !errors.Is(accessAllError, remyErrs.ErrElementNotRegisteredSentinel) {
		err = accessAllError
	}

	// retrieve values from cacheStorage
	return
}

func TryGet[T any](retriever types.DependencyRetriever, tags ...string) (result T) {
	result, _ = Get[T](retriever, tags...)
	return
}

func GetWithPairs[T any](
	retriever types.DependencyRetriever, elements []types.InstancePair[any], tags ...string,
) (result T, err error) {
	subInjector := New(injopts.CacheOptNone, retriever.ReflectOpts(), retriever)
	for _, element := range elements {
		var (
			opts       = injopts.KeyOptsFromStruct(subInjector.ReflectOpts())
			typeSeeker = element.Value
			bindKey    types.BindKey
		)
		if element.InterfaceValue != nil {
			opts |= injopts.KeyOptIgnorePointer
			typeSeeker = element.InterfaceValue
		}
		if bindKey, err = utils.GetElemKey(typeSeeker, opts); err != nil {
			return
		}

		if err = subInjector.BindElem(bindKey, element.Value, types.BindOptions{Tag: element.Tag}); err != nil {
			return
		}
	}

	return Get[T](subInjector, tags...)
}

func GetWith[T any](
	retriever types.DependencyRetriever,
	binder func(injector types.Injector) error, tags ...string,
) (result T, err error) {
	subInjector := New(injopts.CacheOptNone, retriever.ReflectOpts(), retriever)
	if err = binder(subInjector); err != nil {
		return
	}
	return Get[T](subInjector, tags...)
}
