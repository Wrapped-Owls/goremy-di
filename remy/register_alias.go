package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// RegisterAs registers Iface as an explicit alias to an already-registered Impl,
// resolved by a direct lookup instead of DuckTypeElements discovery.
//
// The cast is the identity conversion, and it is what makes the alias
// compile-time checked: `return impl` only compiles when Impl implements Iface.
// The bindFunc selects the alias lifecycle.
//
// Receives: Injector (required); bindFunc (required); cast (required); tag (optional)
func RegisterAs[Iface, Impl any](
	i Injector,
	bindFunc func(types.Binder[Iface]) Bind[Iface],
	cast func(Impl) Iface,
	optTag ...string,
) {
	Register(mustInjector(i), bindFunc(
		func(retriever types.DependencyRetriever) (result Iface, err error) {
			var impl Impl
			if impl, err = Get[Impl](retriever, optTag...); err != nil {
				return result, err
			}
			return cast(impl), nil
		},
	), optTag...)
}
