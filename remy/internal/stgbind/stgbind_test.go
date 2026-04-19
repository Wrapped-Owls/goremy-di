package stgbind

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

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
				t.Errorf("NewStorage returned a storage object that does not implement types.Storage[types.BindKey]")
			}

			if storage == nil {
				t.Fatalf("NewStorage returned nil storage where a %s expected", tc.expectCall)
			}
		})
	}
}
