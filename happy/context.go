package happy

import "iter"

type Context interface {
	GetIndex() any
	GetCurrent() any
	GetData(name string) (any, bool)
	With(index, data any) Context
	WithAll(data ...any) Context
}

type context struct {
	index any
	data  any
	next  *context
}

func NewContext() Context {
	return (*context)(nil)
}

func ContextOf(data ...any) Context {
	return NewContext().WithAll(data...)
}

func (c *context) WithAll(data ...any) Context {
	var ret Context = c
	for i, d := range data {
		ret = ret.With(i, d)
	}
	return ret
}

func (c *context) GetIndex() any {
	return c.index
}

func (c *context) GetCurrent() any {
	return c.data
}

func (c *context) Iter() iter.Seq[Context] {
	return func(yield func(c Context) bool) {
		for curr := c; curr != nil; curr = curr.next {
			if !yield(curr) {
				break
			}
		}
	}
}

func (c *context) GetData(name string) (any, bool) {
	for curr := range c.Iter() {
		data, ok := Get(curr.GetCurrent(), name)
		if ok {
			return data, true
		}
	}
	return nil, false
}

func (c *context) With(index, data any) Context {
	return &context{index, data, c}
}
