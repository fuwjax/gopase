package when

import (
	"reflect"
	"testing"
)

/*
Asserts the actual value is deeply equal to the expected value.
*/
func AssertEqual(t *testing.T, actual, expected any) bool {
	if expected == nil {
		return AssertNil(t, actual)
	} else if actual == nil {
		t.Errorf("actual nil, expected %v", expected)
	} else if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v, expected %v", actual, expected)
	} else {
		return true
	}
	return false
}

/*
Asserts the actual value is nil.
*/
func AssertNil(t *testing.T, actual any) bool {
	if actual != nil {
		t.Errorf("actual %v, expected nil", actual)
		return false
	}
	return true
}

/*
Asserts the actual Error() message matches the expected string.
*/
func AssertError(t *testing.T, actual error, expected string) bool {
	if expected == "" {
		return AssertNil(t, actual)
	} else if actual == nil {
		t.Errorf("actual nil, expected %v", expected)
	} else if actual.Error() != expected {
		t.Errorf("actual %v, expected %v", actual, expected)
	} else {
		return true
	}
	return false
}

func AssertNonZero[T any](t *testing.T, actual T) bool {
	var zero T
	if reflect.DeepEqual(actual, zero) {
		t.Errorf("actual %v, expected non-zero", actual)
		return false
	}
	return true
}

func AssertZero[T any](t *testing.T, actual T) bool {
	var zero T
	if !reflect.DeepEqual(actual, zero) {
		t.Errorf("actual %v, expected zero", actual)
		return false
	}
	return true
}

func AssertFalse(t *testing.T, actual bool) bool {
	if actual {
		t.Errorf("actual %v, expected false", actual)
	}
	return !actual
}

func AssertTrue(t *testing.T, actual bool) bool {
	if !actual {
		t.Errorf("actual %v, expected false", actual)
	}
	return actual
}

func AssertSlices[T ~[]E, E any](t *testing.T, actual, expected T) bool {
	for i, a := range actual {
		if i > len(expected) {
			t.Errorf("actual %v, expected EOF", a)
		}
		if !AssertEqual(t, a, expected[i]) {
			return false
		}
	}
	return AssertEqual(t, len(actual), len(expected))
}
