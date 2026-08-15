package injcontainer

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func Register[T any](ij types.Injector, keyTag string, bind types.Bind[T]) error {
	bindOpts := types.BindOptions{Tag: keyTag}
	return registerNewDep[T](ij, bind, bindOpts)
}

func RegisterWithOverride[T any](ij types.Injector, keyTag string, bind types.Bind[T]) error {
	return registerNewDep[T](ij, bind, types.BindOptions{Tag: keyTag, SoftOverride: true})
}

func registerNewDep[T any](ij types.Injector, bind types.Bind[T], opts types.BindOptions) error {
	elementType := utils.NewKeyElem[T]()
	var (
		value any = bind
		err   error
	)
	if bindType := bind.Type(); bindType == types.BindInstance || bindType == types.BindSingleton {
		// entering here reports a cycle when the chain circles back into this key,
		// instead of failing as an unregistered element
		var retriever types.DependencyRetriever = ij
		if scoped := retriever.RetrieverFor(elementType, opts.Tag); scoped != nil {
			retriever = scoped
		}
		if value, err = bind.Generates(retriever); err != nil {
			return err
		}
	}

	return ij.BindElem(elementType, value, opts)
}
