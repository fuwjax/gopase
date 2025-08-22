package parser

import "fmt"

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
Literal = "'" SingleExpr "'" / '"' DoubleExpr '"'
CharClass = Pattern
Ref = Name

Comment = '#' (!EOL .)*
Name = [_a-zA-Z] [_a-zA-Z0-9]*
Pattern = '[' ('\]' / [^\]])+ ']'
SingleExpr = ("\'" / [^'])*
DoubleExpr = ('\"' / [^"])*
WS = [ \t]*
EOL = [\n\r]
EOF = !.
`

/*
Yet another Peg-grammar parser. Should be identical to Bootstrap.
*/
var Peg = Deferred(func() *Parser {
	g, err := Bootstrap.Parse(grammar)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return g.(*Grammar).Parser("Grammar", PegHandler)
})
