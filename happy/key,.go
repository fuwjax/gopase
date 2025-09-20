package happy

import (
	"strings"

	"github.com/fuwjax/gopase/funki"
)

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
	return SelfKey{}
}

func At() Key {
	return IndexKey{}
}

func Lit(name string) Key {
	return LiteralKey{name}
}

func Dotted(keys ...Key) Key {
	if len(keys) == 1 {
		return keys[0]
	}
	return DottedKey{keys}
}

func Bracket(name string, keys ...Key) Key {
	if len(keys) == 0 {
		keys = nil
	}
	return BracketKey{name, keys}
}

type LiteralKey struct {
	name string
}

func (lit LiteralKey) Resolve(context Context) (any, bool) {
	return lit.name, true
}

func (lit LiteralKey) String() string {
	return lit.name
}

type SelfKey struct{}

func (SelfKey) Resolve(context Context) (any, bool) {
	return context.GetCurrent(), true
}

func (SelfKey) String() string {
	return "."
}

type IndexKey struct{}

func (IndexKey) Resolve(context Context) (any, bool) {
	return context.GetIndex(), true
}

func (IndexKey) String() string {
	return "@"
}

type DottedKey struct {
	brackets []Key
}

func (dot DottedKey) Resolve(context Context) (any, bool) {
	for _, bracket := range dot.brackets {
		data, ok := bracket.Resolve(context)
		if !ok || data == nil {
			return nil, false
		}
		context = context.With(nil, data)
	}
	return context.GetCurrent(), true
}

func (dot DottedKey) String() string {
	return strings.Join(funki.Apply(dot.brackets, Key.String), ".")
}

type BracketKey struct {
	name string
	args []Key
}

func (bkey BracketKey) Resolve(context Context) (any, bool) {
	data, ok := context.GetData(bkey.name)
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

func (bkey BracketKey) String() string {
	return bkey.name + strings.Join(funki.Apply(bkey.args, func(k Key) string { return "[" + k.String() + "]" }), "")
}
