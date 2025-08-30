package mustache

import (
	"iter"
	"strings"

	"github.com/fuwjax/gopase/parser"
)

const mustacheGrammar = `
Template = Line* !.
Line = CommentLine / (Tag / Plain)* EOL
CommentLine = WS '{{' CommentTag '}}' WS EOL
Tag = '{{' (CommentTag / Section / Triple / Ampersand / Ref) '}}'
Section = '#' Ident '}}' Tag '{{/' Ident
Triple = '{' Ident '}'
Ampersand = '&' Ident
Ref = Ident
CommentTag = '!' Comment
Comment = ('}'? [^}])*
Ident = WS Name WS
Name = [._a-zA-Z] [._a-zA-Z0-9]*
Plain = ('{'? [^{\r\n])+
WS = [ \t]*
EOL = '\r\n' / [\r\n] / !.
`

var Grammar = parser.Preserve2(parser.Bootstrap(mustacheGrammar))

func Compile(template string) (*Template, error) {
	grammar, err := Grammar()
	if err != nil {
		return nil, err
	}
	result, err := parser.Parse("Template", grammar, parser.WrapHandler(mustacheHandler{}), template)
	if err != nil {
		return nil, err
	}
	return result.(*Template), nil
}

type mustacheHandler struct{}

func (h mustacheHandler) Template(results iter.Seq2[string, any]) (any, error) {
	renderers := parser.Merge(parser.Cast[[]Renderer](parser.ListOf(results, "Line")))
	return &Template{renderers}, nil
}

func (h mustacheHandler) Line(results iter.Seq2[string, any]) (any, error) {
	result := make([]Renderer, 0)
	var plain strings.Builder
	for name, snippet := range results {
		switch name {
		case "Plain", "EOL":
			plain.WriteString(snippet.(string))
		case "Tag", "CommentLine":
			if plain.Len() > 0 {
				result = append(result, &Plaintext{plain.String()})
				plain.Reset()
			}
			result = append(result, snippet.(Renderer))
		}
	}
	if plain.Len() > 0 {
		result = append(result, &Plaintext{plain.String()})
	}
	return result, nil
}

func (h mustacheHandler) CommentLine(results iter.Seq2[string, any]) (any, error) {
	_, result := parser.FirstOf(results, "CommentTag")
	return result, nil
}

func (h mustacheHandler) Tag(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "CommentTag", "Section", "Ref", "Ampersand", "Triple")
	return value, nil
}

func (h mustacheHandler) Section(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Ident")
	_, tag := parser.FirstOf(results, "Tag")
	return &Section{value.(string), tag.(Renderer)}, nil
}

func (h mustacheHandler) Ref(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Ident")
	return &Reference{value.(string), true}, nil
}

func (h mustacheHandler) Ampersand(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Ident")
	return &Reference{value.(string), false}, nil
}

func (h mustacheHandler) Triple(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Ident")
	return &Reference{value.(string), false}, nil
}

func (h mustacheHandler) CommentTag(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Comment")
	return &Comment{value.(string)}, nil
}

func (h mustacheHandler) Ident(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Name")
	return value.(string), nil
}
