package injcontainer

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stgbind"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func GetWithPairs[T any](
	retriever types.DependencyRetriever, keyTag string,
	newScope types.ScopeFactory, elements ...types.BindEntry,
) (result T, err error) {
	stg := stgbind.NewStorage(injopts.CacheOptNone, uint(len(elements)))
	scope := newScope(retriever, stg)
	for _, element := range elements {
		value, bindKey := element.Entry()
		if bindKey == nil { // Gen a bindKey if none is provided
			err = remyErrs.ErrImpossibleIdentifyType{Type: (*T)(nil)}
			return
		}
		if err = scope.BindElem(
			bindKey,
			value,
			types.BindOptions{Tag: element.Tag()},
		); err != nil {
			return
		}
	}

	return Get[T](scope, keyTag)
}

func GetWith[T any](
	retriever types.DependencyRetriever, keyTag string,
	newScope types.ScopeFactory, binder func(injector types.Injector) error,
) (result T, err error) {
	scope := newScope(retriever, stgbind.NewStorage(injopts.CacheOptNone, 4))
	if err = binder(scope); err != nil {
		return
	}

	return Get[T](scope, keyTag)
}
