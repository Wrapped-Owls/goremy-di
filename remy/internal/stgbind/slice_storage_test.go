package stgbind

import (
	"errors"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func TestSliceStorage_Set(t *testing.T) {
	type setOp struct {
		op           func(*SliceStorage[types.BindKey]) (bool, error)
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
					func(s *SliceStorage[types.BindKey]) (bool, error) { return s.Set(types.KeyElem[uint]{}, 7) },
					false,
				},
				{
					func(s *SliceStorage[types.BindKey]) (bool, error) { return s.Set(types.KeyElem[uint]{}, 11) },
					true,
				},
			},
		},
		{
			name: "Named",
			ops: []setOp{
				{
					func(s *SliceStorage[types.BindKey]) (bool, error) {
						return s.SetNamed(types.KeyElem[string]{}, "lang", "dart")
					},
					false,
				},
				{
					func(s *SliceStorage[types.BindKey]) (bool, error) {
						return s.SetNamed(types.KeyElem[string]{}, "lang", "go")
					},
					true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSliceStorage[types.BindKey](injopts.CacheOptAllowOverride, 2)
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

func TestSliceStorage_Set__Override(t *testing.T) {
	stg := NewSliceStorage[types.BindKey](injopts.CacheOptNone, 2)
	for _, tc := range generateStorageTestCases() {
		t.Run(tc.name, func(t *testing.T) {
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

func TestSliceStorage_Get(t *testing.T) {
	type entry struct {
		key   types.BindKey
		value any
	}
	tests := []struct {
		name    string
		entries []entry
	}{
		{
			name:    "single entry",
			entries: []entry{{types.KeyElem[uint]{}, 42}},
		},
		{
			name: "multiple distinct entries coexist",
			entries: []entry{
				{types.KeyElem[uint]{}, 42},
				{types.KeyElem[string]{}, "hello"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSliceStorage[types.BindKey](injopts.CacheOptNone, uint(len(tt.entries)))

			for _, e := range tt.entries {
				if _, err := stg.Get(e.key); !errors.Is(
					err,
					remyErrs.ErrElementNotRegisteredSentinel,
				) {
					t.Errorf(
						"expected ErrElementNotRegistered before Set for %T, got %v",
						e.key,
						err,
					)
				}
			}
			for _, e := range tt.entries {
				if _, err := stg.Set(e.key, e.value); err != nil {
					t.Fatalf("unexpected error on Set: %v", err)
				}
			}
			for _, e := range tt.entries {
				result, err := stg.Get(e.key)
				if err != nil {
					t.Errorf("unexpected error for %T: %v", e.key, err)
				}
				if result != e.value {
					t.Errorf("wanted %v, got %v", e.value, result)
				}
			}
		})
	}
}

func TestSliceStorage_GetNamed(t *testing.T) {
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
			stg := NewSliceStorage[types.BindKey](injopts.CacheOptNone, 2)

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

func TestSliceStorage_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		opts    injopts.CacheConfOption
		setup   func(*SliceStorage[types.BindKey])
		tag     string
		wantErr error
		wantLen int
	}{
		{
			name:    "ReturnAllDisabled",
			opts:    injopts.CacheOptNone,
			setup:   func(s *SliceStorage[types.BindKey]) { _, _ = s.Set(types.KeyElem[uint]{}, 1) },
			tag:     "",
			wantErr: remyErrs.ErrConfigNotAllowReturnAll,
		},
		{
			name: "unnamed entries returned, named excluded",
			opts: injopts.CacheOptReturnAll,
			setup: func(s *SliceStorage[types.BindKey]) {
				_, _ = s.Set(types.KeyElem[uint]{}, 1)
				_, _ = s.Set(types.KeyElem[string]{}, "a")
				_, _ = s.SetNamed(types.KeyElem[bool]{}, "flag", true)
			},
			tag:     "",
			wantLen: 2,
		},
		{
			name: "named entries returned by tag, unnamed excluded",
			opts: injopts.CacheOptReturnAll,
			setup: func(s *SliceStorage[types.BindKey]) {
				_, _ = s.Set(types.KeyElem[uint]{}, 1)
				_, _ = s.SetNamed(types.KeyElem[string]{}, "lang", "go")
				_, _ = s.SetNamed(types.KeyElem[bool]{}, "lang", true)
			},
			tag:     "lang",
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := NewSliceStorage[types.BindKey](tt.opts, 4)
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
