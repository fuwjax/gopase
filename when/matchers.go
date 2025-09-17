package when

import (
	"iter"
	"testing"
)

func MatchSeq[T any](expected ...T) Matcher[iter.Seq[T]] {
	return func(t *testing.T, actual iter.Seq[T]) {
		var i int
		for elem := range actual {
			AssertEqual(t, elem, expected[i])
			i++
		}
	}
}
