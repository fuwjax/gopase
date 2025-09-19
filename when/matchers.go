package when

import (
	"iter"
	"testing"
)

func MatchEqual[T any](expected T) Matcher[T] {
	return func(t *testing.T, actual T) bool {
		return AssertEqual(t, actual, expected)
	}
}

func MatchSeq[T any](assert Assertion[T], expected ...T) Matcher[iter.Seq[T]] {
	return func(t *testing.T, actual iter.Seq[T]) bool {
		var i int
		for e := range actual {
			if !assert(t, e, expected[i]) {
				return false
			}
			i++
		}
		return true
	}
}

func MatchSlice[T ~[]E, E any](assert Assertion[E], expected T) Matcher[T] {
	return func(t *testing.T, actual T) bool {
		for i, e := range actual {
			if !assert(t, e, expected[i]) {
				return false
			}
		}
		return AssertEqual(t, len(actual), len(expected))
	}
}
