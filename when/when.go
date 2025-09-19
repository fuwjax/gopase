package when

import (
	"testing"
)

type Assertion[T any] func(t *testing.T, actual, expected T) bool

/*
Matcher is a partial Assertion. It encapsulates the logic to determine the expected
result rather than the actual expected instance.
*/
type Matcher[T any] func(t *testing.T, actual T) bool

type Expectation[T any] interface {
	Expect(t *testing.T, expected T) T
	ExpectMatch(t *testing.T, matcher Matcher[T]) T
	ExpectSuccess(t *testing.T) T
	ExpectFailure(t *testing.T) T
	ExpectError(t *testing.T, msg string) T
}

func You[T any](actual T) Expectation[T] {
	return &singleEx[T]{actual}
}

func YouErr[T any](actual T, err error) Expectation[T] {
	return &errorEx[T]{actual, err}
}

func YouOk[T any](actual T, ok bool) Expectation[T] {
	return &boolEx[T]{actual, ok}
}

type OpExpectation[T any] interface {
	Expect(t *testing.T, expected T)
	ExpectMatch(t *testing.T, matcher Matcher[T])
	ExpectSuccess(t *testing.T)
	ExpectFailure(t *testing.T)
	ExpectError(t *testing.T, msg string)
}

type WhenOp[T any] func() T

func YouDo[T any](name string, op WhenOp[T]) OpExpectation[T] {
	return &singleOpEx[T]{name, op}
}

type WhenOpErr[T any] func() (T, error)

func YouDoErr[T any](name string, op WhenOpErr[T]) OpExpectation[T] {
	return &errorOpEx[T]{name, op}
}

type WhenOpOk[T any] func() (T, bool)

func YouDoOk[T any](name string, op WhenOpOk[T]) OpExpectation[T] {
	return &boolOpEx[T]{name, op}
}

type singleEx[T any] struct {
	actual T
}

func (e *singleEx[T]) Expect(t *testing.T, expected T) T {
	AssertEqual(t, e.actual, expected)
	return e.actual
}

func (e *singleEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) T {
	matcher(t, e.actual)
	return e.actual
}

func (e *singleEx[T]) ExpectSuccess(t *testing.T) T {
	AssertNonZero(t, e.actual)
	return e.actual
}

func (e *singleEx[T]) ExpectFailure(t *testing.T) T {
	AssertZero(t, e.actual)
	return e.actual
}

func (e *singleEx[T]) ExpectError(t *testing.T, msg string) T {
	t.Errorf("could not expect error, use YouErr() instead")
	return e.actual
}

type errorEx[T any] struct {
	actual    T
	actualErr error
}

func (e *errorEx[T]) Expect(t *testing.T, expected T) T {
	AssertNil(t, e.actualErr)
	AssertEqual(t, e.actual, expected)
	return e.actual
}

func (e *errorEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) T {
	AssertNil(t, e.actualErr)
	matcher(t, e.actual)
	return e.actual
}

func (e *errorEx[T]) ExpectSuccess(t *testing.T) T {
	AssertNil(t, e.actualErr)
	return e.actual
}

func (e *errorEx[T]) ExpectFailure(t *testing.T) T {
	AssertNonZero(t, e.actualErr)
	return e.actual
}

func (e *errorEx[T]) ExpectError(t *testing.T, msg string) T {
	AssertError(t, e.actualErr, msg)
	return e.actual
}

type boolEx[T any] struct {
	actual   T
	actualOk bool
}

func (e *boolEx[T]) Expect(t *testing.T, expected T) T {
	AssertTrue(t, e.actualOk)
	AssertEqual(t, e.actual, expected)
	return e.actual
}

func (e *boolEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) T {
	AssertTrue(t, e.actualOk)
	matcher(t, e.actual)
	return e.actual
}

func (e *boolEx[T]) ExpectSuccess(t *testing.T) T {
	AssertTrue(t, e.actualOk)
	return e.actual
}

func (e *boolEx[T]) ExpectFailure(t *testing.T) T {
	AssertFalse(t, e.actualOk)
	return e.actual
}

func (e *boolEx[T]) ExpectError(t *testing.T, msg string) T {
	t.Errorf("could not expect error, use YouErr() instead")
	return e.actual
}

type singleOpEx[T any] struct {
	name string
	op   WhenOp[T]
}

func (e *singleOpEx[T]) Expect(t *testing.T, expected T) {
	t.Run(e.name, func(t *testing.T) {
		You(e.op()).Expect(t, expected)
	})
}

func (e *singleOpEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) {
	t.Run(e.name, func(t *testing.T) {
		You(e.op()).ExpectMatch(t, matcher)
	})
}

func (e *singleOpEx[T]) ExpectSuccess(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		You(e.op()).ExpectSuccess(t)
	})
}

func (e *singleOpEx[T]) ExpectFailure(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		You(e.op()).ExpectFailure(t)
	})
}

func (e *singleOpEx[T]) ExpectError(t *testing.T, msg string) {
	t.Errorf("could not expect error, use YouDoErr() instead")
}

type errorOpEx[T any] struct {
	name string
	op   WhenOpErr[T]
}

func (e *errorOpEx[T]) Expect(t *testing.T, expected T) {
	t.Run(e.name, func(t *testing.T) {
		YouErr(e.op()).Expect(t, expected)
	})
}

func (e *errorOpEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) {
	t.Run(e.name, func(t *testing.T) {
		YouErr(e.op()).ExpectMatch(t, matcher)
	})
}

func (e *errorOpEx[T]) ExpectSuccess(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		YouErr(e.op()).ExpectSuccess(t)
	})
}

func (e *errorOpEx[T]) ExpectFailure(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		YouErr(e.op()).ExpectFailure(t)
	})
}

func (e *errorOpEx[T]) ExpectError(t *testing.T, msg string) {
	t.Run(e.name, func(t *testing.T) {
		YouErr(e.op()).ExpectError(t, msg)
	})
}

type boolOpEx[T any] struct {
	name string
	op   WhenOpOk[T]
}

func (e *boolOpEx[T]) Expect(t *testing.T, expected T) {
	t.Run(e.name, func(t *testing.T) {
		YouOk(e.op()).Expect(t, expected)
	})
}

func (e *boolOpEx[T]) ExpectMatch(t *testing.T, matcher Matcher[T]) {
	t.Run(e.name, func(t *testing.T) {
		YouOk(e.op()).ExpectMatch(t, matcher)
	})
}
func (e *boolOpEx[T]) ExpectSuccess(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		YouOk(e.op()).ExpectSuccess(t)
	})
}

func (e *boolOpEx[T]) ExpectFailure(t *testing.T) {
	t.Run(e.name, func(t *testing.T) {
		YouOk(e.op()).ExpectFailure(t)
	})
}

func (e *boolOpEx[T]) ExpectError(t *testing.T, msg string) {
	t.Errorf("could not expect error, use YouDoErr() instead")
}
