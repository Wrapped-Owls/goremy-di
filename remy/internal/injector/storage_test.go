package injector

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func TestElementsStorage_Set(t *testing.T) {
	var checkpoints uint8
	// Checks if panics when trying to override
	defer func() {
		r := recover()
		if r != nil {
			t.Error("Function is panicking")
			t.FailNow()
		}
		if checkpoints != 2 {
			t.Error("Test panic before reaching all checkpoints")
		}
	}()

	stg := NewElementsStorage[string](injopts.CacheOptAllowOverride, types.ReflectionOptions{})
	_ = stg.Set("value", 7)
	checkpoints++
	_ = stg.Set("value", 11)

	_ = stg.SetNamed("tools", "lang", "dart")
	checkpoints++
	_ = stg.SetNamed("tools", "lang", "go")
}

func TestElementsStorage_Set__Override(t *testing.T) {
	testCases := generateStorageTestCases()

	stg := NewElementsStorage[string](injopts.CacheOptNone, types.ReflectionOptions{})
	for _, toTest := range testCases {
		t.Run(
			toTest.name, func(t *testing.T) {
				var checkpoints uint8
				// Checks if panics when trying to override
				defer func() {
					r := recover()
					if r == nil {
						t.Error("Function did not panic")
						t.FailNow()
					}
					if checkpoints != 1 {
						t.Error("Test panic before reaching the first checkpoint")
					}
				}()

				toTest.setterFunc(stg, toTest.values[0])
				checkpoints++
				toTest.setterFunc(stg, toTest.values[1])
			},
		)
	}
}

func generateStorageTestCases() []struct {
	name       string
	values     [2]any
	setterFunc func(types.Storage[string], any)
} {
	return []struct {
		name       string
		values     [2]any
		setterFunc func(types.Storage[string], any)
	}{
		{
			name:   "Set Without Key",
			values: [2]any{7, 11},
			setterFunc: func(stg types.Storage[string], receive any) {
				_ = stg.Set("value", receive)
			},
		},
		{
			name:   "Set Named",
			values: [2]any{"go", "flutter"},
			setterFunc: func(stg types.Storage[string], receive any) {
				_ = stg.SetNamed("value", "tool", receive)
			},
		},
	}
}
