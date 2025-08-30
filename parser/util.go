package parser

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func AssertEqual(t *testing.T, actual, expected any) {
	if actual == nil {
		t.Errorf("actual nil, expected %v", expected)
	} else if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v, expected %v", actual, expected)
	}
}

func AssertNil(t *testing.T, actual any) {
	if actual != nil {
		t.Errorf("actual %v, expected nil", actual)
	}
}

func AssertError(t *testing.T, actual error, expected string) {
	if actual.Error() != expected {
		t.Errorf("actual %v, expected %v", actual, expected)
	}
}

func Preserve2[T, E any](value T, err E) func() (T, E) {
	return func() (T, E) {
		return value, err
	}
}

func Bool2Err[T any](value T, ok bool) (T, error) {
	if ok {
		return value, nil
	}
	return value, fmt.Errorf("value not ok")
}

func Map2Func[K comparable, T any](m map[K]T) func(K) T {
	return func(key K) T {
		return m[key]
	}
}

func Map2Func2[K comparable, T any](m map[K]T) func(K) (T, bool) {
	return func(key K) (T, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func FirstOf[K comparable, T any](i iter.Seq2[K, T], keys ...K) (K, T) {
	for k, t := range i {
		if slices.Index(keys, k) >= 0 {
			return k, t
		}
	}
	var k K
	var t T
	return k, t
}

func ListOf[K comparable, T any](i iter.Seq2[K, T], key K) []T {
	results := make([]T, 0)
	for k, t := range i {
		if k == key {
			results = append(results, t)
		}
	}
	return results
}

func Merge[T any](tss [][]T) []T {
	result := make([]T, 0)
	for _, ts := range tss {
		result = append(result, ts...)
	}
	return result
}

func MapOf[T any, K comparable](source []T, key func(T) K) map[K]T {
	result := make(map[K]T, len(source))
	for _, t := range source {
		result[key(t)] = t
	}
	return result
}

func Apply[F, T any](source []F, xform func(F) T) []T {
	result := make([]T, len(source))
	for i, value := range source {
		result[i] = xform(value)
	}
	return result
}

func Filter[T any](source []T, pred func(T) bool) []T {
	result := make([]T, 0, len(source))
	for _, value := range source {
		if pred(value) {
			result = append(result, value)
		}
	}
	return result
}

func ToString() func(fmt.Stringer) string {
	return fmt.Stringer.String
}

func Cast[T any](source []any) []T {
	result := make([]T, len(source))
	for i, value := range source {
		if value != nil {
			result[i] = value.(T)
		}
	}
	return result
}

type PolyError struct {
	Errors []error
}

func (es *PolyError) Add(e error) {
	es.Errors = append(es.Errors, e)
}

func (es *PolyError) Error() string {
	return strings.Join(Apply(es.Errors, error.Error), "\n")
}

func Deferred[T any](initializer func() T) func() T {
	var cache T
	var inited bool
	return func() T {
		if !inited {
			cache = initializer()
			inited = true
		}
		return cache
	}
}
