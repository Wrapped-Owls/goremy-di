package injcontainer

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stdinj"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// StandardScope is the ScopeFactory used when nothing decorates the retriever.
func StandardScope(
	parent types.DependencyRetriever, storage types.Storage[types.BindKey],
) types.Injector {
	return stdinj.NewWithStorage(stdinj.Options{}, storage, parent)
}

func TryGet[T any](retriever types.DependencyRetriever, keyTag string) (result T) {
	result, _ = Get[T](retriever, keyTag)
	return
}
