package funki

import (
	"iter"
	"slices"
)

/*
Returns the first key-value pair from the iterator that matches one of the keys.
Returns the zero-zero pair if nothing matches.
*/
func FirstOf[K comparable, T any](i iter.Seq2[K, T], keys ...K) (K, T) {
	return First(FilterKeys(i, keys...))
}

func ListOf[T any](i iter.Seq2[string, any], keys ...string) []T {
	return slices.Collect(Cast[T](Values(FilterKeys(i, keys...))))
}

func First[K any, T any](i iter.Seq2[K, T]) (K, T) {
	for k, t := range i {
		return k, t
	}
	var k K
	var t T
	return k, t
}

func FilterKeys[K comparable, T any](i iter.Seq2[K, T], keys ...K) iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		i(func(k K, t T) bool {
			if slices.Index(keys, k) >= 0 {
				return yield(k, t)
			}
			return true
		})
	}
}

func FilterNonNil[T comparable](i iter.Seq[T]) iter.Seq[T] {
	var zero T
	return func(yield func(T) bool) {
		i(func(t T) bool {
			if t == zero {
				return true
			}
			return yield(t)
		})
	}
}

func Values[K, T any](i iter.Seq2[K, T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		i(func(k K, t T) bool {
			return yield(t)
		})
	}
}

func Cast[T any](i iter.Seq[any]) iter.Seq[T] {
	var zero T
	return func(yield func(T) bool) {
		i(func(v any) bool {
			if v == nil {
				return yield(zero)
			}
			return yield(v.(T))
		})
	}
}
