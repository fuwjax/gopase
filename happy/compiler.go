package happy

import (
	"fmt"
	"iter"
	"strings"
	"sync"

	"github.com/fuwjax/gopase/funki"
	"github.com/fuwjax/gopase/parser"
)

// (^-^) should be the index/map key/field name in a section
const happyGrammar = `
Template = Content EOF
Content = (Plain? Tag)* Plain?
Tag = Open (Comment / Value / Include / Section / Else / Partial / Override) Close

Comment = '#' Text?
Value = Key
Include = '>' KeyName
Section = '*' Key Close Content Open '/' Text?
Else = '!' Key Close Content Open '/' Text?
Partial = '=' KeyName Close Content Open '/' Text?
Override = '>>' KeyName Close Content Open '/' Text?

Text = (!Close .)+
KeyName = Key / Name
Name = WS '"' String '"' WS
String = (Escape / Chs)*
Chs = [^"]+
Escape = '""'
Key = WS (Bracket ('.' Bracket)* / Dot / At) WS
Dot = '.'
At = '@'
Bracket = Ident ('[' KeyName ']')*
Ident = [_a-zA-Z0-9]+
Plain = (!Open .)+
Open = '(^' / WS '( ^'
Close = '^)' / '^ )' WS
WS = [ \t\r\n]*
EOF = !.
`

// ParserFrom exists exclusively for testing rules directly
var ParserFrom = sync.OnceValue(func() parser.ParserFrom {
	return parser.NewParserFrom(happyGrammar, happyHandler{})
})
var templateParser = sync.OnceValue(func() parser.Parser[Template] {
	return parser.NewParser[Template]("Template", happyGrammar, happyHandler{})
})

func Compile(template string) (Template, error) {
	return templateParser()(template)
}

type happyHandler struct{}

func (h happyHandler) Template(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Content")
	return value.(Template), nil
}

func (h happyHandler) Content(results iter.Seq2[string, any]) (any, error) {
	result := make([]Template, 0)
	for name, snippet := range results {
		switch name {
		case "Plain":
			result = append(result, Plaintext(snippet.(string)))
		case "Tag":
			result = append(result, snippet.(Template))
		}
	}
	return Content(result), nil
}

func (h happyHandler) Tag(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Value", "Else", "Include", "Section", "Partial", "Override", "Otherwise")
	return value, nil
}

func (h happyHandler) Otherwise(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Text")
	return Plaintext(value.(string)), nil
}

func (h happyHandler) Value(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Key")
	return Reference(value.(Key)), nil
}

func (h happyHandler) Include(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "KeyName")
	return Include(value.(Key)), nil
}

func (h happyHandler) Section(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Key")
	_, content := funki.FirstOf(results, "Content")
	return Section(value.(Key), content.(Template)), nil
}

func (h happyHandler) Else(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Key")
	_, content := funki.FirstOf(results, "Content")
	return Invert(value.(Key), content.(Template)), nil
}

func (h happyHandler) Partial(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "KeyName")
	_, content := funki.FirstOf(results, "Content")
	return Partial(value.(Key), content.(Template)), nil
}

func (h happyHandler) Override(results iter.Seq2[string, any]) (any, error) {
	//	_, value := parser.FirstOf(results, "KeyName")
	//	_, content := parser.FirstOf(results, "Content")
	//	return &Override{value.(Key), content.(Template)}, nil
	return nil, fmt.Errorf("not implemented")
}

func (h happyHandler) KeyName(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Key", "Name")
	return value, nil
}

func (h happyHandler) Name(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "String")
	return Lit(value.(string)), nil
}

func (h happyHandler) String(results iter.Seq2[string, any]) (any, error) {
	var sb strings.Builder
	for _, text := range results {
		sb.WriteString(text.(string))
	}
	return sb.String(), nil
}

func (happyHandler) Escape(results iter.Seq2[string, any]) (any, error) {
	return "\"", nil
}

func (happyHandler) Key(results iter.Seq2[string, any]) (any, error) {
	name, _ := funki.FirstOf(results, "Dot", "At")
	switch name {
	case "At":
		return At(), nil
	case "Dot":
		return Dot(), nil
	}
	brackets := funki.ListOf[Key](results, "Bracket")
	return Dotted(brackets...), nil
}

func (happyHandler) Bracket(results iter.Seq2[string, any]) (any, error) {
	_, name := funki.FirstOf(results, "Ident")
	args := funki.ListOf[Key](results, "KeyName")
	return Bracket(name.(string), args...), nil
}
