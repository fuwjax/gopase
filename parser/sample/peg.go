package sample

import (
	"fmt"

	"github.com/fuwjax/gopase/parser"
)

const grammar = `
Grammar = Line (EOL Line)* EOL? EOF
Line = Rule / Comment / WS
Rule = WS Name WS '=' WS Expr WS
Expr = Seq (WS '/' WS Seq)*
Seq = Prefix (WS Prefix)*
Prefix = AndExpr / NotExpr / Suffix
AndExpr = '&' WS Suffix
NotExpr = '!' WS Suffix
Suffix = OptExpr / RepExpr / ReqExpr / Primary
OptExpr = Primary WS '?'
RepExpr = Primary WS '*'
ReqExpr = Primary WS '+'
Primary = Dot / ParExpr / Literal / CharClass / Ref
Dot = '.'
ParExpr = '(' WS Expr WS ')'
Literal = SingleLit / DoubleLit
CharClass = Pattern
Ref = Name

Comment = '#' (!EOL .)*
Name = [_a-zA-Z] [_a-zA-Z0-9]*
Pattern = '[' ("\\]" / [^\]])+ ']'
SingleLit = "'" ("\\" SingleEscape / SinglePlain)* "'"
DoubleLit = '"' ("\\" DoubleEscape / DoublePlain)* '"'
SingleEscape = [\\'nrt]
DoubleEscape = [\\"nrt] 
SinglePlain = [^\\']+
DoublePlain = [^\\"]+
WS = [ \t]*
EOL = [\n\r]
EOF = !.
`

/*
Yet another Peg-grammar parser. Should be identical to Bootstrap.
*/
var Peg = func() *parser.Grammar {
	grammar, err := parser.Bootstrap(grammar)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return grammar
}()
