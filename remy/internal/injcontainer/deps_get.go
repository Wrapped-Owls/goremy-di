package injcontainer

import (
	"errors"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func Get[T any](retriever types.DependencyRetriever, keyTag string) (element T, err error) {
	elementType := utils.NewKeyElem[T]()

	if scoped := retriever.RetrieverFor(elementType, keyTag); scoped != nil {
		retriever = scoped
	}

	var bind any
	// search in dynamic injections that needed to run a given function
	if bind, err = retriever.RetrieveBind(elementType, keyTag); err == nil {
		if typedBind, assertOk := bind.(types.Bind[T]); assertOk {
			return typedBind.Generates(retriever)
		}
		if instanceBind, assertOk := utils.Satisfies[T](bind); assertOk {
			return instanceBind, nil
		}
		err = remyErrs.ErrTypeCastInRuntime{ActualValue: bind, Expected: (*T)(nil)}
	}

	if !utils.IsInterface[T]() {
		return // If the element is indeed an interface, we skip the guess recover
	}

	// MultiBinding alone must not pay this O(n) scan on an interface miss
	if opts, ok := types.ResolveOptionsOf(retriever); ok &&
		!opts.Is(injopts.ResolveOptDuckTyping) {
		return
	}

	// Start to search for every element if it is configured in this way
	foundElement, accessAllError := getByGuess[T](retriever, keyTag)
	if accessAllError == nil {
		element = foundElement
		err = nil
	} else if !errors.Is(accessAllError, remyErrs.ErrConfigNotAllowReturnAll) &&
		!shouldIgnoreGuessError(accessAllError, elementType) {
		err = accessAllError
	}

	// retrieve values from cacheStorage
	return
}

func GetAll[T any](
	retriever types.DependencyRetriever, keyTag string,
) (resultList []T, err error) {
	var elementList []any
	if elementList, err = retriever.GetAll(keyTag); err != nil {
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
		bindKey := utils.NewKeyElem[T]()
		err = remyErrs.ErrElementNotRegistered{Key: bindKey}
	}

	return
}

func shouldIgnoreGuessError(checkErr error, requestedKey types.BindKey) bool {
	notRegistered, ok := remyErrs.CheckError[remyErrs.ErrElementNotRegistered](checkErr)
	var missingKey types.BindKey
	missingKey, ok = notRegistered.Key.(types.BindKey)
	return ok && missingKey.ID() == requestedKey.ID()
}

func checkSavedAsBind[T any](
	retriever types.DependencyRetriever, checkElem any,
) (foundElem *T, err error) {
	if genericBind, assertOk := utils.Satisfies[types.GuessableBind](checkElem); assertOk {
		// Check if the returned value can implement the requested interface
		if _, ok := utils.Satisfies[T](genericBind.DefaultValue()); !ok {
			// Check again but now for pointer value, this fallback works because is known that
			// the registered value cannot be an interface at this point
			if _, ok = utils.Satisfies[T](genericBind.PointerValue()); !ok {
				// Value is not applicable to the type
				return nil, nil
			}
		}

		var anyVal any
		if anyVal, err = genericBind.GenAsAny(retriever); err != nil {
			return
		} else if bindElem, ok := utils.Satisfies[T](anyVal); ok {
			foundElem = &bindElem
		}
	}
	return
}

func getByGuess[T any](
	retriever types.DependencyRetriever, keyTag string,
) (element T, err error) {
	var elementList []T
	if elementList, err = GetAll[T](retriever, keyTag); err != nil {
		return element, err
	}

	totalFound := len(elementList)
	if totalFound == 1 {
		element = elementList[0]
		return element, nil
	}

	bindKey := utils.NewKeyElem[T]()
	err = remyErrs.ErrMultipleDIDuckTypingCandidates{Type: bindKey, Count: totalFound}
	if totalFound == 0 {
		err = remyErrs.ErrElementNotRegistered{Key: bindKey}
	}

	return element, err
}
