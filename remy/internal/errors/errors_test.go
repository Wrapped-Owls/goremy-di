package errors

import (
	"errors"
	"testing"
)

type mockTestError struct {
	Message string
}

func (e mockTestError) Error() string {
	return e.Message
}

type wrapperTestError struct {
	msg string
	err error
}

func (e wrapperTestError) Error() string {
	return e.msg
}

func (e wrapperTestError) Unwrap() error {
	return e.err
}

func TestBaseErrorCheckerIs_Comprehensive(t *testing.T) {
	type customTestErrorChecker struct {
		baseErrorChecker[customTestErrorChecker, *customTestErrorChecker]
		mockTestError
	}

	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "same raw error type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "hello world"}},
			target: customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "same pointer error type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "hello world"}},
			target: &customTestErrorChecker{},
			want:   true,
		},
		{
			name: "same error as pointer type",
			err: customTestErrorChecker{
				mockTestError: mockTestError{Message: "not nil-pointer"},
			},
			target: &customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "check pointer generated with raw type",
			err:    customTestErrorChecker{mockTestError: mockTestError{Message: "pointer"}},
			target: customTestErrorChecker{},
			want:   true,
		},
		{
			name:   "different error type",
			err:    customTestErrorChecker{},
			target: errors.New("some other error"),
			want:   false,
		},
		{
			name:   "nil error",
			err:    nil,
			target: customTestErrorChecker{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.err, tt.target); got != tt.want {
				t.Errorf("errors.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckError(t *testing.T) {
	mockTarget := mockTestError{Message: "target_error_type"}

	tests := []struct {
		name      string
		checkErr  any
		wantOk    bool
		wantFound mockTestError
	}{
		{
			name:      "direct match T",
			checkErr:  mockTarget,
			wantOk:    true,
			wantFound: mockTarget,
		},
		{
			name:      "direct match *T (non nil)",
			checkErr:  &mockTarget,
			wantOk:    true,
			wantFound: mockTarget,
		},
		{
			name:      "direct match *T (nil)",
			checkErr:  (*mockTestError)(nil),
			wantOk:    false,
			wantFound: mockTestError{},
		},
		{
			name:      "no match different type",
			checkErr:  errors.New("random message"),
			wantOk:    false,
			wantFound: mockTestError{},
		},
		{
			name:      "nil check error",
			checkErr:  nil,
			wantOk:    false,
			wantFound: mockTestError{},
		},
		{
			name:      "check error is not an error (string)",
			checkErr:  "not an error",
			wantOk:    false,
			wantFound: mockTestError{},
		},
		{
			name: "unwrap leads to match T",
			checkErr: wrapperTestError{
				msg: "outer",
				err: mockTarget,
			},
			wantOk:    true,
			wantFound: mockTarget,
		},
		{
			name: "unwrap leads to match *T (non nil)",
			checkErr: wrapperTestError{
				msg: "outer",
				err: &mockTarget,
			},
			wantOk:    true,
			wantFound: mockTarget,
		},
		{
			name: "unwrap chain matches T",
			checkErr: wrapperTestError{
				msg: "outer",
				err: wrapperTestError{
					msg: "middle",
					err: mockTarget,
				},
			},
			wantOk:    true,
			wantFound: mockTarget,
		},
		{
			name: "unwrap chain fails match",
			checkErr: wrapperTestError{
				msg: "outer",
				err: wrapperTestError{
					msg: "middle",
					err: errors.New("wrong type"),
				},
			},
			wantOk:    false,
			wantFound: mockTestError{},
		},
		{
			name: "unwrap chain leads to nil check",
			checkErr: wrapperTestError{
				msg: "outer",
				err: nil,
			},
			wantOk:    false,
			wantFound: mockTestError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundErr, ok := CheckError[mockTestError](tt.checkErr)
			if ok != tt.wantOk {
				t.Errorf("CheckError() ok = %v, want %v", ok, tt.wantOk)
			}
			if foundErr.Message != tt.wantFound.Message {
				t.Errorf(
					"CheckError() found error message = %v, want %v",
					foundErr.Message, tt.wantFound.Message,
				)
			}
		})
	}
}
