package stgbind

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// NewStorage returns the most efficient Storage implementation for the given
// expected element count:
//
//   - size == 1: SingleStorage (flat fields, zero heap overhead)
//   - 1 < size <= 4: SliceStorage (linear scan, one contiguous alloc)
//   - size > 4: ElementsStorage (hash-map, O(1) lookup)
func NewStorage(opts injopts.CacheConfOption, size uint) types.Storage[types.BindKey] {
	switch {
	case size == 1:
		return NewSingleStorage(opts)
	case size <= 4:
		return NewSliceStorage(opts, size)
	default:
		return NewElementsStorage[types.BindKey](opts)
	}
}
