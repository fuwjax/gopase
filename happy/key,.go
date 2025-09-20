package happy

import (
	"strings"

	"github.com/fuwjax/gopase/funki"
)

// Key and the following functions are public for testing. These are not intended for use outside the package.

type Key interface {
	Resolve(context Context) (any, bool)
	String() string
}

func ResolveName(key Key, context Context) string {
	data, ok := key.Resolve(context)
	if !ok || data == nil {
		return key.String()
	}
	name, ok := String(data)
	if ok {
		return name
	}
	return ""
}

func Dot() Key {
	return &selfKey{}
}

func At() Key {
	return &indexKey{}
}

func Lit(name string) Key {
	return &literalKey{name}
}

func Dotted(keys ...Key) Key {
	if len(keys) == 1 {
		return keys[0]
	}
	return &dottedKey{keys}
}

func Bracket(name string, keys ...Key) Key {
	if len(keys) == 0 {
		keys = nil
	}
	return &bracketKey{name, keys}
}

type literalKey struct {
	name string
}

func (lit *literalKey) Resolve(context Context) (any, bool) {
	return lit.name, true
}

func (lit *literalKey) String() string {
	return lit.name
}

type selfKey struct{}

func (*selfKey) Resolve(context Context) (any, bool) {
	return context.GetData(), true
}

func (*selfKey) String() string {
	return "."
}

type indexKey struct{}

func (*indexKey) Resolve(context Context) (any, bool) {
	return context.GetIndex(), true
}

func (*indexKey) String() string {
	return "@"
}

type dottedKey struct {
	brackets []Key
}

func (dot *dottedKey) Resolve(context Context) (any, bool) {
	for _, bracket := range dot.brackets {
		data, ok := bracket.Resolve(context)
		if !ok || data == nil {
			return nil, false
		}
		context = ContextOf(nil, data)
	}
	return context.GetData(), true
}

func (dot *dottedKey) String() string {
	return strings.Join(funki.Apply(dot.brackets, Key.String), ".")
}

type bracketKey struct {
	name string
	args []Key
}

func (bkey *bracketKey) Resolve(context Context) (any, bool) {
	data, ok := context.Get(bkey.name)
	if !ok {
		return nil, false
	}
	args := make([]any, len(bkey.args))
	for i, arg := range bkey.args {
		args[i], ok = arg.Resolve(context)
		if !ok || args[i] == nil {
			args[i] = arg.String()
		}
	}
	return Call(data, args)
}

func (bkey *bracketKey) String() string {
	return bkey.name + strings.Join(funki.Apply(bkey.args, func(k Key) string { return "[" + k.String() + "]" }), "")
}
