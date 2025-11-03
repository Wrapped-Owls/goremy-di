package errors

import (
	"testing"
)

func TestErrTypeCastInRuntime_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ErrTypeCastInRuntime
		want string
	}{
		{
			name: "nil actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: nil,
				Expected:    "string",
			},
			want: "unable to find/cast the element with given type: string",
		},
		{
			name: "string actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: "hello",
				Expected:    42,
			},
			want: "unable to cast `string` to given type `int`",
		},
		{
			name: "int actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: 42,
				Expected:    "string",
			},
			want: "unable to cast `int` to given type `string`",
		},
		{
			name: "float actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: 3.14,
				Expected:    "string",
			},
			want: "unable to cast `float64` to given type `string`",
		},
		{
			name: "slice actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: []int{1, 2, 3},
				Expected:    "string",
			},
			want: "unable to cast `[]int` to given type `string`",
		},
		{
			name: "map actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: map[string]int{"key": 1},
				Expected:    "string",
			},
			want: "unable to cast `map[string]int` to given type `string`",
		},
		{
			name: "struct actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: struct{ Name string }{Name: "John"},
				Expected:    "string",
			},
			want: "unable to cast `struct { Name string }` to given type `string`",
		},
		{
			name: "interface actual value",
			err: ErrTypeCastInRuntime{
				ActualValue: interface{}(42),
				Expected:    "string",
			},
			want: "unable to cast `int` to given type `string`",
		},
		{
			name: "both nil",
			err: ErrTypeCastInRuntime{
				ActualValue: nil,
				Expected:    nil,
			},
			want: "unable to find/cast the element with given type: <nil>",
		},
		{
			name: "both same type",
			err: ErrTypeCastInRuntime{
				ActualValue: "hello",
				Expected:    "world",
			},
			want: "unable to cast `string` to given type `string`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("ErrTypeCastInRuntime.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
