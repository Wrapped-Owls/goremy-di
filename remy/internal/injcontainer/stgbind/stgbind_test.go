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
			name: "Slice Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewSliceStorage[types.BindKey](injopts.CacheOptNone, length)
			},
			sizes: []uint{1, 2, 3, 4},
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
			name: "Slice Element Storage",
			constructor: func(length uint) types.Storage[types.BindKey] {
				return NewSliceStorage[types.BindKey](injopts.CacheOptNone, length)
			},
			sizes: []uint{1, 2, 3, 4},
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

func TestStorage_ForEach(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		storage types.Storage[types.BindKey]
		entries map[string]any
	}{
		{
			name:    "single storage visits its only entry",
			storage: NewSingleStorage[types.BindKey](injopts.CacheOptNone),
			entries: map[string]any{"": 42},
		},
		{
			name:    "slice storage visits tagged and untagged entries",
			storage: NewSliceStorage[types.BindKey](injopts.CacheOptNone, 3),
			entries: map[string]any{"": 1, "second": 2, "third": 3},
		},
		{
			name:    "map storage visits tagged and untagged entries",
			storage: NewElementsStorage[types.BindKey](injopts.CacheOptNone),
			entries: map[string]any{"": 1, "second": 2, "third": 3},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			keys := map[string]types.BindKey{
				"":       types.KeyElem[int]{},
				"second": types.KeyElem[string]{},
				"third":  types.KeyElem[bool]{},
			}
			for tag, value := range testCase.entries {
				var err error
				if tag == "" {
					_, err = testCase.storage.Set(keys[tag], value)
				} else {
					_, err = testCase.storage.SetNamed(keys[tag], tag, value)
				}
				if err != nil {
					t.Fatalf("set %q: %v", tag, err)
				}
			}

			seen := map[string]any{}
			testCase.storage.ForEach(func(tag string, value any) bool {
				seen[tag] = value
				return true
			})

			if len(seen) != len(testCase.entries) {
				t.Fatalf("visited %v, want %v", seen, testCase.entries)
			}
			for tag, want := range testCase.entries {
				if seen[tag] != want {
					t.Errorf("entry %q = %v, want %v", tag, seen[tag], want)
				}
			}
		})
	}
}

func TestStorage_ForEachStopsEarly(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		storage types.Storage[types.BindKey]
	}{
		{name: "slice storage", storage: NewSliceStorage[types.BindKey](injopts.CacheOptNone, 3)},
		{name: "map storage", storage: NewElementsStorage[types.BindKey](injopts.CacheOptNone)},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			for index, key := range []types.BindKey{
				types.KeyElem[int]{}, types.KeyElem[string]{}, types.KeyElem[bool]{},
			} {
				if _, err := testCase.storage.SetNamed(
					key,
					string(rune('a'+index)),
					index,
				); err != nil {
					t.Fatalf("set: %v", err)
				}
			}

			var visited int
			testCase.storage.ForEach(func(string, any) bool {
				visited++
				return false // asking to stop after the first entry
			})

			if visited != 1 {
				t.Fatalf("visited %d entries, want to stop at 1", visited)
			}
		})
	}
}
