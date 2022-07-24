package utils

import (
	"reflect"
	"testing"
)

type (
	empty interface{}
	a     interface{ A() string }
	b     interface{ B() bool }
	c     interface{ C() }
	d     interface{ D() float64 }
	full  interface {
		a
		b
		c
		d
	}
)

type TypeTestCase struct {
	expectedStr string
	reflectType reflect.Type
}

var typeElements = [...]TypeTestCase{
	{expectedStr: "interface {  }", reflectType: reflect.TypeOf((*empty)(nil)).Elem()},
	{expectedStr: "interface { A func() string }", reflectType: reflect.TypeOf((*a)(nil)).Elem()},
	{expectedStr: "interface { B func() bool }", reflectType: reflect.TypeOf((*b)(nil)).Elem()},
	{expectedStr: "interface { C func() }", reflectType: reflect.TypeOf((*c)(nil)).Elem()},
	{expectedStr: "interface { D func() float64 }", reflectType: reflect.TypeOf((*d)(nil)).Elem()},
	{
		expectedStr: "interface { A func() string; B func() bool; C func(); D func() float64 }",
		reflectType: reflect.TypeOf((*full)(nil)).Elem(),
	},
}

func TestBuildDuckInterfaceType(t *testing.T) {
	for _, typeCase := range typeElements {
		result := buildDuckInterfaceType(typeCase.reflectType)
		if result != typeCase.expectedStr {
			t.Errorf("Duck typing was not correct: Expected `%s`\nReceived: `%s`", typeCase.expectedStr, result)
		}
	}
}

func extractType[T any](fromElement bool) (result reflect.Type) {
	if fromElement {
		var element T
		result, _ = GetElemType[T](element)
		return
	}
	result, _ = GetType[T]()
	return
}

func runDuckTestCases(t *testing.T, testCases [6]TypeTestCase) {
	for _, testCase := range testCases {
		result := buildDuckInterfaceType(testCase.reflectType.Elem())
		if testCase.expectedStr != result {
			t.Errorf(
				"Type obtained is not found correctly. Expected: `%s`\nReceived: `%s`",
				testCase.expectedStr, result,
			)
		}
	}
}

func TestGetType(t *testing.T) {
	/* Build test cases when getting type using generics alongside reflection */
	var testCases = [...]TypeTestCase{
		{expectedStr: typeElements[0].expectedStr, reflectType: extractType[empty](false)},
		{expectedStr: typeElements[1].expectedStr, reflectType: extractType[a](false)},
		{expectedStr: typeElements[2].expectedStr, reflectType: extractType[b](false)},
		{expectedStr: typeElements[3].expectedStr, reflectType: extractType[c](false)},
		{expectedStr: typeElements[4].expectedStr, reflectType: extractType[d](false)},
		{expectedStr: typeElements[5].expectedStr, reflectType: extractType[full](false)},
	}

	runDuckTestCases(t, testCases)

	/* Start to test getting type from the element */
	testCases = [...]TypeTestCase{
		{expectedStr: typeElements[0].expectedStr, reflectType: extractType[empty](true)},
		{expectedStr: typeElements[1].expectedStr, reflectType: extractType[a](true)},
		{expectedStr: typeElements[2].expectedStr, reflectType: extractType[b](true)},
		{expectedStr: typeElements[3].expectedStr, reflectType: extractType[c](true)},
		{expectedStr: typeElements[4].expectedStr, reflectType: extractType[d](true)},
		{expectedStr: typeElements[5].expectedStr, reflectType: extractType[full](true)},
	}

	runDuckTestCases(t, testCases)
}
