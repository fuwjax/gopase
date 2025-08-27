package sample

import (
	"fmt"
	"iter"
	"strings"

	"github.com/fuwjax/gopase/parser"
)

const csvGrammar = `
Records = Record (EOL Record)* EOL? EOF
Record = Field (',' Field)* 
Field = Quoted / Bare
Quoted = WS '"' Inner '"' WS
Inner = ([^"] / '""')*
Bare = [^,\n\r]*
WS = [ \t]*
EOL = [\n\r]
EOF = !.
`

/*
A CSV parser to illustrate some fundamentals
*/
var CSV = func() parser.Parser[[][]string] {
	grammar, err := parser.Bootstrap(csvGrammar)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return parser.NewParser[[][]string]("Grammar", grammar, csvHandler{})
}()

type csvHandler struct{}

func (h csvHandler) Records(results iter.Seq2[string, any]) (any, error) {
	return parser.ListOf(results, "Record"), nil
}

func (h csvHandler) Record(results iter.Seq2[string, any]) (any, error) {
	return parser.ListOf(results, "Field"), nil
}

func (h csvHandler) Field(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Quoted", "Bare")
	return value, nil
}

func (h csvHandler) Quoted(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Inner")
	value = strings.ReplaceAll(value.(string), "\"\"", "\"")
	return value, nil
}
