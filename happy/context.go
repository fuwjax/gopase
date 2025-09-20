package happy

import "iter"

/*
Context manages the search space for Key lookups. The default implementation manages this
space as an immutable stack of index-data pairs. Both GetIndex and GetData retrieve their respective elements
from the top of the stack. Get searches the stack starting at the top for any element containing
the given name. With returns a new stack with the given pair added to the top of the current stack.

It is not required, when implementing your own Context, to mimic this behavior.
*/
type Context interface {
	/*
		Returns the index currently tracked by this context.

		The default implementation returns the index at the top of the stack.
	*/
	GetIndex() any
	/*
		Returns the data currently tracked by this context.

		The default implementation returns the data at the top of the stack.
	*/
	GetData() any
	/*
		Returns a property located within this context. Ok will the value was successfully found.

		The default implementation searches the stack for the first property named name.
	*/
	Get(name string) (value any, ok bool)
	/*
		Creates a new mashup of the current context and the given index-data pair.

		The default implementation pushes an index-data pair onto the stack, and returns the new stack.
	*/
	With(index, data any) Context
}

type context struct {
	index any
	data  any
	next  *context
}

func ContextOf(data ...any) Context {
	return WithAll((*context)(nil), data...)
}

func WithAll(c Context, data ...any) Context {
	for i, d := range data {
		c = c.With(i, d)
	}
	return c
}

func (c *context) GetIndex() any {
	return c.index
}

func (c *context) GetData() any {
	return c.data
}

func (c *context) seq() iter.Seq[Context] {
	return func(yield func(c Context) bool) {
		for curr := c; curr != nil; curr = curr.next {
			if !yield(curr) {
				break
			}
		}
	}
}

func (c *context) Get(name string) (any, bool) {
	for curr := range c.seq() {
		data, ok := Get(curr.GetData(), name)
		if ok {
			return data, true
		}
	}
	return nil, false
}

func (c *context) With(index, data any) Context {
	return &context{index, data, c}
}
