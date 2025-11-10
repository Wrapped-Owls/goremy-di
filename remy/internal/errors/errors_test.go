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
