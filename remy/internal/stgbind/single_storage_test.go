package stgbind

import (
	"errors"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func TestSingleStorage_Set(t *testing.T) {
	type setOp struct {
		op           func(*SingleStorage[types.BindKey]) (bool, error)
		wantOverride bool
	}
	tests := []struct {
		name string
		ops  []setOp
	}{
		{
			name: "Unnamed",
			ops: []setOp{
				{
					func(s *SingleStorage[types.BindKey]) (bool, error) { return s.Set(types.KeyElem[uint]{}, 7) },
					false,
				},
				{
					func(s *SingleStorage[types.BindKey]) (bool, error) { return s.Set(types.KeyElem[uint]{}, 11) },
					true,
				},
			},
		},
		{
			name: "Named",
			ops: []setOp{
				{
					func(s *SingleStorage[types.BindKey]) (bool, error) {
						return s.SetNamed(types.KeyElem[string]{}, "lang", "dart")
					},
					false,
				},
				{
					func(s *SingleStorage[types.BindKey]) (bool, error) {
						return s.SetNamed(types.KeyElem[string]{}, "lang", "go")
					},
					true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSingleStorage[types.BindKey](injopts.CacheOptAllowOverride)
			for _, op := range tt.ops {
				wasOverridden, err := op.op(stg)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if wasOverridden != op.wantOverride {
					t.Errorf("wanted overridden %v, got %v", op.wantOverride, wasOverridden)
				}
			}
		})
	}
}

func TestSingleStorage_Set__Override(t *testing.T) {
	// Each test case needs its own SingleStorage because only one entry fits.
	for _, tc := range generateStorageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
			stg := NewSingleStorage[types.BindKey](injopts.CacheOptNone)

			wasOverridden, err := tc.setterFunc(stg, tc.values[0])
			if err != nil {
				t.Errorf("unexpected error on first set: %v", err)
			}
			if wasOverridden {
				t.Errorf("unexpected override on first set of %v", tc.values[0])
			}

			wasOverridden, err = tc.setterFunc(stg, tc.values[1])
			if err == nil {
				t.Error("expected error when overriding without permission, got none")
			}
			if wasOverridden {
				t.Errorf("unexpected override for %v", tc.values[1])
			}
		})
	}
}

func TestSingleStorage_Set__SecondDistinctEntry(t *testing.T) {
	stg := NewSingleStorage[types.BindKey](injopts.CacheOptNone)
	if _, err := stg.Set(types.KeyElem[uint]{}, 7); err != nil {
		t.Fatalf("unexpected error on first set: %v", err)
	}

	tests := []struct {
		name string
		op   func() (bool, error)
	}{
		{
			"different key",
			func() (bool, error) { return stg.Set(types.KeyElem[string]{}, "second") },
		},
		{
			"same key different name",
			func() (bool, error) { return stg.SetNamed(types.KeyElem[uint]{}, "named", 99) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.op(); !errors.Is(err, remyErrs.ErrAlreadyBoundSentinel) {
				t.Errorf("expected ErrAlreadyBound, got %v", err)
			}
		})
	}
}

func TestSingleStorage_Get(t *testing.T) {
	tests := []struct {
		name    string
		key     types.BindKey
		value   any
		missKey types.BindKey
	}{
		{
			name:    "uint value",
			key:     types.KeyElem[uint]{},
			value:   42,
			missKey: types.KeyElem[string]{},
		},
		{
			name:    "string value",
			key:     types.KeyElem[string]{},
			value:   "hello",
			missKey: types.KeyElem[uint]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSingleStorage[types.BindKey](injopts.CacheOptNone)

			if _, err := stg.Get(tt.key); !errors.Is(
				err,
				remyErrs.ErrElementNotRegisteredSentinel,
			) {
				t.Errorf("expected ErrElementNotRegistered before Set, got %v", err)
			}
			if _, err := stg.Set(tt.key, tt.value); err != nil {
				t.Fatalf("unexpected error on Set: %v", err)
			}
			result, err := stg.Get(tt.key)
			if err != nil {
				t.Errorf("unexpected error after Set: %v", err)
			}
			if result != tt.value {
				t.Errorf("wanted %v, got %v", tt.value, result)
			}
			if _, err = stg.Get(tt.missKey); !errors.Is(
				err,
				remyErrs.ErrElementNotRegisteredSentinel,
			) {
				t.Errorf("expected ErrElementNotRegistered for different key, got %v", err)
			}
		})
	}
}

func TestSingleStorage_GetNamed(t *testing.T) {
	tests := []struct {
		name  string
		key   types.BindKey
		tag   string
		value any
	}{
		{name: "string with lang tag", key: types.KeyElem[string]{}, tag: "lang", value: "go"},
		{name: "int with version tag", key: types.KeyElem[int]{}, tag: "version", value: 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSingleStorage[types.BindKey](injopts.CacheOptNone)

			if _, err := stg.SetNamed(tt.key, tt.tag, tt.value); err != nil {
				t.Fatalf("unexpected error on SetNamed: %v", err)
			}
			result, err := stg.GetNamed(tt.key, tt.tag)
			if err != nil {
				t.Errorf("unexpected error on GetNamed: %v", err)
			}
			if result != tt.value {
				t.Errorf("wanted %v, got %v", tt.value, result)
			}
			if _, err = stg.GetNamed(tt.key, "other"); !errors.Is(
				err,
				remyErrs.ErrElementNotRegisteredSentinel,
			) {
				t.Errorf("expected ErrElementNotRegistered for wrong tag, got %v", err)
			}
			if _, err = stg.Get(tt.key); !errors.Is(err, remyErrs.ErrElementNotRegisteredSentinel) {
				t.Errorf(
					"expected ErrElementNotRegistered for unnamed lookup on named entry, got %v",
					err,
				)
			}
		})
	}
}

func TestSingleStorage_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		opts    injopts.CacheConfOption
		setup   func(*SingleStorage[types.BindKey])
		tag     string
		wantErr error
		wantLen int
	}{
		{
			name:    "ReturnAllDisabled",
			opts:    injopts.CacheOptNone,
			setup:   func(s *SingleStorage[types.BindKey]) { _, _ = s.Set(types.KeyElem[uint]{}, 1) },
			tag:     "",
			wantErr: remyErrs.ErrConfigNotAllowReturnAll,
		},
		{
			name:    "unnamed entry returned",
			opts:    injopts.CacheOptReturnAll,
			setup:   func(s *SingleStorage[types.BindKey]) { _, _ = s.Set(types.KeyElem[uint]{}, 99) },
			tag:     "",
			wantLen: 1,
		},
		{
			name:    "named entry returned by tag",
			opts:    injopts.CacheOptReturnAll,
			setup:   func(s *SingleStorage[types.BindKey]) { _, _ = s.SetNamed(types.KeyElem[string]{}, "lang", "go") },
			tag:     "lang",
			wantLen: 1,
		},
		{
			name:    "tag mismatch returns empty",
			opts:    injopts.CacheOptReturnAll,
			setup:   func(s *SingleStorage[types.BindKey]) { _, _ = s.SetNamed(types.KeyElem[string]{}, "lang", "go") },
			tag:     "other",
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSingleStorage[types.BindKey](tt.opts)
			tt.setup(stg)

			all, err := stg.GetAll(tt.tag)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted error %v, got %v", tt.wantErr, err)
			}
			if len(all) != tt.wantLen {
				t.Errorf("wanted %d entries, got %d", tt.wantLen, len(all))
			}
		})
	}
}
