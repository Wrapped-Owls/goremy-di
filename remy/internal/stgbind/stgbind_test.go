package stgbind

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func TestElementsStorage_Set(t *testing.T) {
	stg := NewElementsStorage[types.BindKey](injopts.CacheOptAllowOverride)
	var (
		wasOverridden bool
		err           error
	)

	checkFunc := func(expectedOverride bool, expectedErr error) {
		if wasOverridden != expectedOverride {
			t.Errorf("Wanted overridden %v, got %v", expectedOverride, wasOverridden)
		}

		if !errors.Is(err, expectedErr) {
			t.Errorf("Wanted error %v, got %v", expectedErr, expectedErr)
		}
	}
	wasOverridden, err = stg.Set(types.KeyElem[uint]{}, 7)
	checkFunc(false, nil)

	wasOverridden, err = stg.Set(types.KeyElem[uint]{}, 11)
	checkFunc(true, nil)

	wasOverridden, err = stg.SetNamed(types.KeyElem[string]{}, "lang", "dart")
	checkFunc(false, nil)
	wasOverridden, err = stg.SetNamed(types.KeyElem[string]{}, "lang", "go")
	checkFunc(true, nil)
}

func TestElementsStorage_Set__Override(t *testing.T) {
	testCases := generateStorageTestCases()

	stg := NewElementsStorage[types.BindKey](injopts.CacheOptNone)
	for _, toTest := range testCases {
		t.Run(
			toTest.name, func(t *testing.T) {
				wasOverridden, err := toTest.setterFunc(stg, toTest.values[0])
				if err != nil {
					t.Errorf("Unexpected error on first set: %v", err)
				}
				if wasOverridden {
					t.Errorf("Unexpected overridden %v", toTest.values[0])
				}

				wasOverridden, err = toTest.setterFunc(stg, toTest.values[1])
				if err == nil {
					t.Error(
						"Expected error when trying to override without override permission, but got none",
					)
				}
				if wasOverridden {
					t.Errorf("Unexpected overridden %v", toTest.values[0])
				}
			},
		)
	}
}

func generateStorageTestCases() []struct {
	name       string
	values     [2]any
	setterFunc func(types.Storage[types.BindKey], any) (bool, error)
} {
	return []struct {
		name       string
		values     [2]any
		setterFunc func(types.Storage[types.BindKey], any) (bool, error)
	}{
		{
			name:   "Set Without Tag",
			values: [2]any{7, 11},
			setterFunc: func(stg types.Storage[types.BindKey], receive any) (bool, error) {
				return stg.Set(types.KeyElem[any]{}, receive)
			},
		},
		{
			name:   "Set Named",
			values: [2]any{"go", "flutter"},
			setterFunc: func(stg types.Storage[types.BindKey], receive any) (bool, error) {
				return stg.SetNamed(types.KeyElem[any]{}, "tool", receive)
			},
		},
	}
}
