package testi

import (
	"reflect"
	"testing"
)

/*
Asserts the actual value is deeply equal to the expected value.
*/
func AssertEqual(t *testing.T, actual, expected any) {
	if actual == nil {
		t.Errorf("actual nil, expected %v", expected)
	} else if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v, expected %v", actual, expected)
	}
}

/*
Asserts the actual value is nil.
*/
func AssertNil(t *testing.T, actual any) {
	if actual != nil {
		t.Errorf("actual %v, expected nil", actual)
	}
}

/*
Asserts the actual Error() message matches the expected string.
*/
func AssertError(t *testing.T, actual error, expected string) {
	if actual.Error() != expected {
		t.Errorf("actual %v, expected %v", actual, expected)
	}
}
