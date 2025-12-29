package stgbind

import (
	"errors"
	"reflect"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"
)

func TestStorageBackend_Set_Get(t *testing.T) {
	tests := []struct {
		name  string
		items []struct {
			key   types.BindKey
			value any
		}
		allowOverride bool
		getKey        types.BindKey
		wantGetValue  any
		wantGetError  bool
		wantSize      int
	}{
		{
			name: "new key with Optimadin Prime",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{
					key:   types.KeyElem[*fixtures.OptimadinPrime]{},
					value: fixtures.NewOptimadinPrime(),
				},
			},
			getKey:       types.KeyElem[*fixtures.OptimadinPrime]{},
			wantGetValue: fixtures.NewOptimadinPrime(),
			wantSize:     1,
		},
		{
			name: "new key with Bumblelf",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{key: types.KeyElem[*fixtures.Bumblelf]{}, value: fixtures.NewBumblelf()},
			},
			getKey:       types.KeyElem[*fixtures.Bumblelf]{},
			wantGetValue: fixtures.NewBumblelf(),
			wantSize:     1,
		},
		{
			name: "allow override - Optimadin to upgraded Class",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{
					key: types.KeyElem[*fixtures.OptimadinPrime]{},
					value: &fixtures.OptimadinPrime{
						DnDBase: fixtures.DnDBase{ClassName: "Original Paladin Commander"},
					},
				},
				{
					key:   types.KeyElem[*fixtures.OptimadinPrime]{},
					value: fixtures.NewOptimadinPrime(),
				},
			},
			allowOverride: true,
			getKey:        types.KeyElem[*fixtures.OptimadinPrime]{},
			wantGetValue:  fixtures.NewOptimadinPrime(),
			wantSize:      1,
		},
		{
			name: "disallow override - Bumblelf stays Bumblelf",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{
					key: types.KeyElem[*fixtures.Bumblelf]{},
					value: &fixtures.Bumblelf{
						DnDBase: fixtures.DnDBase{ClassName: "Original Elf Scout"},
					},
				},
				{key: types.KeyElem[*fixtures.Bumblelf]{}, value: fixtures.NewBumblelf()},
			},
			getKey: types.KeyElem[*fixtures.Bumblelf]{},
			wantGetValue: &fixtures.Bumblelf{
				DnDBase: fixtures.DnDBase{ClassName: "Original Elf Scout"},
			},
			wantSize: 1,
		},
		{
			name: "existing key - Soundbard",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{
					key: types.KeyElem[*fixtures.Soundbard]{},
					value: &fixtures.Soundbard{
						DnDBase: fixtures.DnDBase{ClassName: "Original Waterwave"},
					},
				},
			},
			getKey: types.KeyElem[*fixtures.Soundbard]{},
			wantGetValue: &fixtures.Soundbard{
				DnDBase: fixtures.DnDBase{ClassName: "Original Waterwave"},
			},
			wantSize: 1,
		},
		{
			name: "non-existing key",
			items: []struct {
				key   types.BindKey
				value any
			}{},
			getKey:       types.KeyElem[*fixtures.Ratcheric]{},
			wantGetValue: nil,
			wantGetError: true,
			wantSize:     0,
		},
		{
			name: "multiple items with duplicates - allowOverride false",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{
					key:   types.KeyElem[*fixtures.OptimadinPrime]{},
					value: fixtures.NewOptimadinPrime(),
				},
				{key: types.KeyElem[*fixtures.Bumblelf]{}, value: fixtures.NewBumblelf()},
				{key: types.KeyElem[*fixtures.Ironknight]{}, value: fixtures.NewIronknight()},
				{key: types.KeyElem[*fixtures.Ratcheric]{}, value: fixtures.NewRatcheric()},
				{key: types.KeyElem[*fixtures.Jazogue]{}, value: fixtures.NewJazogue()},
				{key: types.KeyElem[*fixtures.Wheelificer]{}, value: fixtures.NewWheelificer()},
				{
					key:   types.KeyElem[*fixtures.Ironknight]{}, // Duplicate key
					value: fixtures.Ironknight{},
				},
				{
					key:   types.KeyElem[*fixtures.Ratcheric]{}, // Duplicate key
					value: fixtures.Ratcheric{},
				},
			},
			allowOverride: false,
			getKey:        types.KeyElem[*fixtures.Ironknight]{},
			wantGetValue:  fixtures.NewIronknight(),
			wantSize:      6, // Should have 6 unique bots (Ironknight and Ratcheric are duplicates)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newBackend[types.BindKey](10)

			// Set all items
			for i, item := range tt.items {
				allowOverride := tt.allowOverride
				if i == 0 {
					allowOverride = false // First set never allows override
				}
				backend.Set(item.key, item.value, allowOverride)
			}

			if backend.Size() != tt.wantSize {
				t.Errorf("Expected size %d, got %d", tt.wantSize, backend.Size())
			}

			// Get the value
			retrieved, err := backend.Get(tt.getKey)
			if tt.wantGetError {
				if err == nil {
					t.Error("Expected error for non-existing key")
				}
				var expectedErr remyErrs.ErrElementNotRegistered
				if !errors.As(err, &expectedErr) {
					t.Errorf("Expected ErrElementNotRegistered, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				// Use reflect.DeepEqual for comparison
				if !reflect.DeepEqual(retrieved, tt.wantGetValue) {
					t.Errorf("Expected value %v, got %v", tt.wantGetValue, retrieved)
				}
			}
		})
	}
}

func TestStorageBackend_Size(t *testing.T) {
	tests := []struct {
		name  string
		items []struct {
			key   types.BindKey
			value any
		}
		wantSize int
	}{
		{
			name: "empty backend",
			items: []struct {
				key   types.BindKey
				value any
			}{},
			wantSize: 0,
		},
		{
			name: "single item - Jazz",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{types.KeyElem[*fixtures.Jazogue]{}, fixtures.NewJazogue()},
			},
			wantSize: 1,
		},
		{
			name: "multiple items - Autobot team",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{types.KeyElem[*fixtures.OptimadinPrime]{}, fixtures.NewOptimadinPrime()},
				{types.KeyElem[*fixtures.Bumblelf]{}, fixtures.NewBumblelf()},
				{types.KeyElem[*fixtures.Ironknight]{}, fixtures.NewIronknight()},
				{types.KeyElem[*fixtures.Ratcheric]{}, fixtures.NewRatcheric()},
				{types.KeyElem[*fixtures.Jazogue]{}, fixtures.NewJazogue()},
			},
			wantSize: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newBackend[types.BindKey](10)
			for _, item := range tt.items {
				backend.Set(item.key, item.value, false)
			}

			if backend.Size() != tt.wantSize {
				t.Errorf("Expected size %d, got %d", tt.wantSize, backend.Size())
			}
		})
	}
}

func TestStorageBackend_GetAll(t *testing.T) {
	tests := []struct {
		name  string
		items []struct {
			key   types.BindKey
			value any
		}
		wantCount int
	}{
		{
			name: "empty backend",
			items: []struct {
				key   types.BindKey
				value any
			}{},
			wantCount: 0,
		},
		{
			name: "single item - Wheelificer",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{types.KeyElem[*fixtures.Wheelificer]{}, fixtures.NewWheelificer()},
			},
			wantCount: 1,
		},
		{
			name: "multiple items - Decepticon command",
			items: []struct {
				key   types.BindKey
				value any
			}{
				{types.KeyElem[*fixtures.MegadwarfTron]{}, fixtures.NewMegadwarfTron()},
				{types.KeyElem[*fixtures.Sorcerscream]{}, fixtures.NewSorcerscream()},
				{types.KeyElem[*fixtures.Soundbard]{}, fixtures.NewSoundbard()},
				{types.KeyElem[*fixtures.Shocklock]{}, fixtures.NewShocklock()},
			},
			wantCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newBackend[types.BindKey](10)
			for _, item := range tt.items {
				backend.Set(item.key, item.value, false)
			}

			all := backend.GetAll()
			if len(all) != tt.wantCount {
				t.Fatalf("Expected %d items, got %d", tt.wantCount, len(all))
			}

			// Check that all values are present (order may vary)
			for _, item := range tt.items {
				found := false
				for _, v := range all {
					if reflect.DeepEqual(v, item.value) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected value %v not found in GetAll()", item.value)
				}
			}
		})
	}
}
