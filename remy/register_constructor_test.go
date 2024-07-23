package remy

import (
	"math"
	"strconv"
	"testing"
)

func TestRegisterConstructor(t *testing.T) {
	const registerKey = "random_K3Y-register"
	var tests = []struct {
		name              string
		injectionRegister func(injector Injector, calledTimes *uint16)
	}{
		{
			name: "No parameters",
			injectionRegister: func(inj Injector, calledTimes *uint16) {
				constructor := func() string {
					*calledTimes++
					return "test-value"
				}
				RegisterConstructor(inj, Factory[string], constructor, registerKey)
			},
		},
		{
			name: "One Parameter",
			injectionRegister: func(inj Injector, calledTimes *uint16) {
				constructor := func(arg1 string) string {
					*calledTimes++
					return "Hello, " + arg1
				}
				RegisterConstructorArgs1(inj, Factory[string], constructor, registerKey)
				RegisterInstance(inj, "Remy Dependency Injector")
			},
		},
		{
			name: "Two Parameters",
			injectionRegister: func(inj Injector, calledTimes *uint16) {
				constructor := func(arg1 string, arg2 int) string {
					*calledTimes++
					return arg1 + " is " + strconv.Itoa(arg2)
				}
				RegisterInstance(inj, "Universe Answer")
				RegisterInstance(inj, 42)
				RegisterConstructorArgs2(inj, Factory[string], constructor, registerKey)
			},
		},
		{
			name: "Three Parameters",
			injectionRegister: func(inj Injector, calledTimes *uint16) {
				constructor := func(arg1 string, arg2 int, arg3 bool) string {
					*calledTimes++
					return "(" + strconv.FormatBool(arg3) + ").__" + arg1 + " is " + strconv.Itoa(arg2)
				}
				RegisterInstance(inj, "Universe Answer")
				RegisterInstance(inj, 42)
				RegisterInstance(inj, true)
				RegisterConstructorArgs3(inj, Factory[string], constructor, registerKey)
			},
		},
		{
			name: "Four Parameters",
			injectionRegister: func(inj Injector, calledTimes *uint16) {
				constructor := func(arg1 string, arg2 int, arg3 bool, arg4 float32) string {
					*calledTimes++
					result := "(" + strconv.FormatBool(arg3) + ").__" + arg1 + " is " + strconv.Itoa(arg2)
					result += "\tVALUE: " + strconv.FormatFloat(float64(arg4), 'f', -1, 32)
					return result
				}
				RegisterInstance(inj, "Universe Answer")
				RegisterInstance(inj, 42)
				RegisterInstance(inj, true)
				RegisterInstance(inj, float32(math.Pi))
				RegisterConstructorArgs4(inj, Factory[string], constructor, registerKey)
			},
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			inj := NewInjector(Config{CanOverride: false})
			var calledTimes uint16
			tCase.injectionRegister(inj, &calledTimes)

			value, err := DoGet[string](inj, registerKey)
			if err != nil {
				t.Fatal(err)
			}

			if calledTimes <= 0 || value == "" {
				t.Errorf("Constructor was not called, received value `%s`", value)
			}
		})
	}
}
