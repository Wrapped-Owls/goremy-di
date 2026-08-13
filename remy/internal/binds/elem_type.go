package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

// elemType carries the views that depend only on T, so every bind gets them by
// embedding instead of repeating them. Being zero-sized it costs no memory,
// as long as it is embedded before the other fields.
type elemType[T any] struct{}

func (elemType[T]) PointerValue() any {
	return (*T)(nil)
}

func (elemType[T]) DefaultValue() any {
	var defaultValue T
	return defaultValue
}

func (elemType[T]) ElementKey() types.BindKey {
	return types.KeyElem[T]{}
}
