package funki

import (
	"iter"
	"sync"

	"github.com/fuwjax/gopase/parser"
)

const grammar = ""

// ParserFrom exists exclusively for testing rules directly
var ParserFrom = sync.OnceValue(func() parser.ParserFrom {
	return parser.NewParserFrom(grammar, handler{})
})
var parserFunk = sync.OnceValue(func() parser.Parser[Library] {
	return parser.NewParser[Library]("Template", grammar, handler{})
})

func Compile(funkis string) (Library, error) {
	return parserFunk()(funkis)
}

type handler struct{}

func (h handler) Library(results iter.Seq2[string, any]) (any, error) {
	value := ListOf[Function](results, "Declaration")
	return NewLibrary(value), nil
}

func (h handler) Declaration(results iter.Seq2[string, any]) (any, error) {
	_, name := FirstOf(results, "Name")
	_, relation := FirstOf(results, "Relation")
	return NewFunction(name.(string), relation.(Relation)), nil
}

func (h handler) Relation(results iter.Seq2[string, any]) (any, error) {
	_, pattern := FirstOf(results, "Pattern")
	_, production := FirstOf(results, "Production")
	_, relation := FirstOf(results, "Relation")
	if relation == nil {
		return NewRelation(pattern.(Pattern), production.(Relation)), nil
	}
	return NewRelation(pattern.(Pattern), relation.(Relation)), nil
}

func (h handler) Pattern(results iter.Seq2[string, any]) (any, error) {
	name, value := First(results)
	switch name {
	case "Underbar":
		return Underbar{}, nil
	case "Name":
		return NewReference(value.(string)), nil
	case "Number", "Literal":
		return NewExact(value), nil
	default:
		return value, nil
	}
}

func (h handler) ListPattern(results iter.Seq2[string, any]) (any, error) {
	value := ListOf[Pattern](results, "Pattern")
	return NewListPattern(value)
}
