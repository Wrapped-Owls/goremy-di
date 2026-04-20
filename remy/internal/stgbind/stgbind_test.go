package stgbind

import (
	"strconv"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// distinctKeys is a pool of unique keys used by benchmark tables.
// Each element has a different underlying type, so every ID() is unique.
var distinctKeys = []types.BindKey{
	types.KeyElem[uint]{},
	types.KeyElem[string]{},
	types.KeyElem[bool]{},
	types.KeyElem[int]{},
	types.KeyElem[float32]{},
	types.KeyElem[float64]{},
	types.KeyElem[uint8]{},
	types.KeyElem[int8]{},
	types.KeyElem[uint16]{},
	types.KeyElem[int16]{},
}

func BenchmarkStorage_Set(b *testing.B) {
	cases := []struct {
		name        string
		constructor func(length uint) types.Storage[types.BindKey]
		sizes       []uint
	}{
		{
			name: "Single Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewSingleStorage[types.BindKey](injopts.CacheOptNone)
			},
			sizes: []uint{1},
		},
		{
			name:        "Slice Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] { return NewSliceStorage(injopts.CacheOptNone, length) },
			sizes:       []uint{1, 2, 3, 4},
		},
		{
			name: "Map Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewElementsStorage[types.BindKey](injopts.CacheOptNone)
			},
			sizes: []uint{1, 2, 3, 4, 5, 10},
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for _, size := range tc.sizes {
				b.Run(strconv.FormatUint(uint64(size), 10), func(b *testing.B) {
					keys := distinctKeys[:size]
					b.Helper()
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						stg := tc.constructor(size)
						for _, k := range keys {
							_, _ = stg.Set(k, struct{}{})
						}
					}
				})
			}
		})
	}
}

func BenchmarkStorage_Get(b *testing.B) {
	cases := []struct {
		name        string
		constructor func(length uint) types.Storage[types.BindKey]
		sizes       []uint
	}{
		{
			name: "Single Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewSingleStorage[types.BindKey](injopts.CacheOptNone)
			},
			sizes: []uint{1},
		},
		{
			name:        "Slice Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] { return NewSliceStorage(injopts.CacheOptNone, length) },
			sizes:       []uint{1, 2, 3, 4},
		},
		{
			name: "Map Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewElementsStorage[types.BindKey](injopts.CacheOptNone)
			},
			sizes: []uint{1, 2, 3, 4, 5, 10},
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for _, size := range tc.sizes {
				b.Run(strconv.FormatUint(uint64(size), 10), func(b *testing.B) {
					keys := distinctKeys[:size]
					stg := tc.constructor(size)
					for _, k := range keys {
						_, _ = stg.Set(k, struct{}{})
					}

					b.Helper()
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for _, k := range keys {
							_, _ = stg.Get(k)
						}
					}
				})
			}
		})
	}
}

func TestNewStorage(t *testing.T) {
	const opts = injopts.CacheOptNone

	testCases := []struct {
		name       string
		opts       injopts.CacheConfOption
		size       uint
		expectCall string // Used for documentation/clarity, not actual assertion
	}{
		// --- Case: Single Storage (size == 1) ---
		{
			name:       "SizeOneMinimalOpts",
			opts:       opts,
			size:       1,
			expectCall: "NewSingleStorage",
		},
		// --- Case: Slice Storage (size <= 4 and size != 1) ---
		{
			name:       "SizeZero", // Corner case: size 0
			opts:       opts,
			size:       0,
			expectCall: "NewSliceStorage",
		},
		{
			name:       "SizeTwo",
			opts:       opts,
			size:       2,
			expectCall: "NewSliceStorage",
		},
		{
			name:       "SizeFourMaxSlice", // Boundary case: size 4
			opts:       injopts.CacheConfOption(3),
			size:       4,
			expectCall: "NewSliceStorage",
		},
		// --- Case: Elements Storage (default or size > 4) ---
		{
			name:       "SizeFiveMinElements", // Boundary case: size 5
			opts:       opts,
			size:       5,
			expectCall: "NewElementsStorage",
		},
		{
			name:       "LargeSizeElements",
			opts:       opts,
			size:       100,
			expectCall: "NewElementsStorage",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storage := NewStorage(tc.opts, tc.size)

			// Assert the returned storage satisfies the interface
			if _, ok := storage.(types.Storage[types.BindKey]); !ok {
				t.Errorf(
					"NewStorage returned a storage object that does not implement types.Storage[types.BindKey]",
				)
			}

			if storage == nil {
				t.Fatalf("NewStorage returned nil storage where a %s expected", tc.expectCall)
			}
		})
	}
}
