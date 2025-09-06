package happy

import (
	"fmt"
	"strings"

	"github.com/fuwjax/gopase/funki"
)

type Key interface {
	Resolve(context *Context) (any, bool)
	String() string
}

func ResolveName(key Key, context *Context) string {
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

func Invert(key Key) Key {
	return InvertKey{key}
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

type InvertKey struct {
	key Key
}

func (invert InvertKey) Resolve(context *Context) (any, bool) {
	result, ok := invert.key.Resolve(context)
	if !ok || !Truthy(result) {
		return nil, true
	}
	return nil, false
}

func (invert InvertKey) String() string {
	return fmt.Sprintf("!%s", invert.key)
}

type LiteralKey struct {
	name string
}

func (lit LiteralKey) Resolve(context *Context) (any, bool) {
	return lit.name, true
}

func (lit LiteralKey) String() string {
	return lit.name
}

type SelfKey struct{}

func (SelfKey) Resolve(context *Context) (any, bool) {
	return context.Data, true
}

func (SelfKey) String() string {
	return "."
}

type IndexKey struct{}

func (IndexKey) Resolve(context *Context) (any, bool) {
	return context.Index, true
}

func (IndexKey) String() string {
	return "@"
}

type DottedKey struct {
	brackets []Key
}

func (dot DottedKey) Resolve(context *Context) (any, bool) {
	for _, bracket := range dot.brackets {
		data, ok := bracket.Resolve(context)
		if !ok || data == nil {
			return nil, false
		}
		context = context.With(nil, data)
	}
	return context.Data, true
}

func (dot DottedKey) String() string {
	return strings.Join(funki.Apply(dot.brackets, Key.String), ".")
}

type BracketKey struct {
	name string
	args []Key
}

func (bkey BracketKey) Resolve(context *Context) (any, bool) {
	for curr := context; curr != nil; curr = curr.Next {
		data, ok := Get(curr.Data, bkey.name)
		if ok {
			args := make([]any, len(bkey.args))
			for i, arg := range bkey.args {
				args[i], ok = arg.Resolve(context)
				if !ok || args[i] == nil {
					args[i] = arg.String()
				}
			}
			return Call(data, args)
		}
	}
	return nil, false
}

func (bkey BracketKey) String() string {
	return bkey.name + strings.Join(funki.Apply(bkey.args, func(k Key) string { return "[" + k.String() + "]" }), "")
}
