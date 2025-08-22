package parser

/*
The Peg-grammar parser. Bootstrap.Parse() returns a *Grammar.
*/
var Bootstrap = bootstrap()

func bootstrap() *Parser {
	peg := NewGrammar()
	peg.AddRule("Grammar", Seq(Ref("Line"), Rep(Seq(Ref("EOL"), Ref("Line"))), Opt(Ref("EOL")), Ref("EOF")))
	peg.AddRule("Line", Alt(Ref("Rule"), Ref("Comment"), Ref("WS")))
	peg.AddRule("Rule", Seq(Ref("WS"), Ref("Name"), Ref("WS"), Lit("="), Ref("WS"), Ref("Expr"), Ref("WS")))
	peg.AddRule("Expr", Seq(Ref("Seq"), Rep(Seq(Ref("WS"), Lit("/"), Ref("WS"), Ref("Seq")))))
	peg.AddRule("Seq", Seq(Ref("Prefix"), Rep(Seq(Ref("WS"), Ref("Prefix")))))
	peg.AddRule("Prefix", Alt(Ref("AndExpr"), Ref("NotExpr"), Ref("Suffix")))
	peg.AddRule("AndExpr", Seq(Lit("&"), Ref("WS"), Ref("Suffix")))
	peg.AddRule("NotExpr", Seq(Lit("!"), Ref("WS"), Ref("Suffix")))
	peg.AddRule("Suffix", Alt(Ref("OptExpr"), Ref("RepExpr"), Ref("ReqExpr"), Ref("Primary")))
	peg.AddRule("OptExpr", Seq(Ref("Primary"), Ref("WS"), Lit("?")))
	peg.AddRule("RepExpr", Seq(Ref("Primary"), Ref("WS"), Lit("*")))
	peg.AddRule("ReqExpr", Seq(Ref("Primary"), Ref("WS"), Lit("+")))
	peg.AddRule("Primary", Alt(Ref("Dot"), Ref("ParExpr"), Ref("Literal"), Ref("CharClass"), Ref("Ref")))
	peg.AddRule("Dot", Lit("."))
	peg.AddRule("ParExpr", Seq(Lit("("), Ref("WS"), Ref("Expr"), Ref("WS"), Lit(")")))
	peg.AddRule("Literal", Alt(Seq(Lit("'"), Ref("SingleExpr"), Lit("'")), Seq(Lit("\""), Ref("DoubleExpr"), Lit("\""))))
	peg.AddRule("CharClass", Ref("Pattern"))
	peg.AddRule("Ref", Ref("Name"))

	peg.AddRule("Comment", Seq(Lit("#"), Rep(Seq(Not(Ref("EOL")), Dot()))))
	peg.AddRule("Name", Seq(Cls("[_a-zA-Z]"), Rep(Cls("[_a-zA-Z0-9]"))))
	peg.AddRule("Pattern", Seq(Lit("["), Req(Alt(Lit("\\]"), Cls("[^\\]]"))), Lit("]")))
	peg.AddRule("SingleExpr", Rep(Alt(Lit("\\'"), Cls("[^']"))))
	peg.AddRule("DoubleExpr", Rep(Alt(Lit("\\\""), Cls("[^\"]"))))
	peg.AddRule("WS", Rep(Cls("[ \\t]")))
	peg.AddRule("EOL", Cls("[\\n\\r]"))
	peg.AddRule("EOF", Not(Dot()))
	return NewParser("Grammar", PegHandler, peg)
}

var PegHandler = WrapHandler(pegHandler{})

type pegHandler struct{}

func (p pegHandler) Grammar(result *ParseResult) (any, error) {
	rules := Cast[*Rule](Filter(ListOf(result.Results(), "Line"), func(r any) bool { return r != nil }))
	grammar := NewGrammar()
	for _, rule := range rules {
		grammar.Add(rule)
	}
	return grammar, nil
}

func (p pegHandler) Line(result *ParseResult) (any, error) {
	_, rule := FirstOf(result.Results(), "Rule")
	return rule, nil
}

func (p pegHandler) Rule(result *ParseResult) (any, error) {
	_, name := FirstOf(result.Results(), "Name")
	_, expr := FirstOf(result.Results(), "Expr")
	return NewRule(name.(string), expr.(Expr)), nil
}

func (p pegHandler) Expr(result *ParseResult) (any, error) {
	seqs := Cast[Expr](ListOf(result.Results(), "Seq"))
	return Alt(seqs...), nil
}

func (p pegHandler) Seq(result *ParseResult) (any, error) {
	prefixs := Cast[Expr](ListOf(result.Results(), "Prefix"))
	return Seq(prefixs...), nil
}

func (p pegHandler) Prefix(result *ParseResult) (any, error) {
	_, value := FirstOf(result.Results(), "AndExpr", "NotExpr", "Suffix")
	return value, nil
}

func (p pegHandler) AndExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Suffix")
	return See(expr.(Expr)), nil
}

func (p pegHandler) NotExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Suffix")
	return Not(expr.(Expr)), nil
}

func (p pegHandler) Suffix(result *ParseResult) (any, error) {
	_, value := FirstOf(result.Results(), "OptExpr", "RepExpr", "ReqExpr", "Primary")
	return value, nil
}

func (p pegHandler) OptExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Primary")
	return Opt(expr.(Expr)), nil
}

func (p pegHandler) RepExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Primary")
	return Rep(expr.(Expr)), nil
}

func (p pegHandler) ReqExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Primary")
	return Req(expr.(Expr)), nil
}

func (p pegHandler) Primary(result *ParseResult) (any, error) {
	_, value := FirstOf(result.Results(), "Dot", "ParExpr", "Literal", "CharClass", "Ref")
	return value, nil
}

func (p pegHandler) Dot(result *ParseResult) (any, error) {
	return Dot(), nil
}

func (p pegHandler) ParExpr(result *ParseResult) (any, error) {
	_, expr := FirstOf(result.Results(), "Expr")
	return expr, nil
}

func (p pegHandler) Literal(result *ParseResult) (any, error) {
	_, value := FirstOf(result.Results(), "SingleExpr", "DoubleExpr")
	return Lit(value.(string)), nil
}

func (p pegHandler) CharClass(result *ParseResult) (any, error) {
	_, pattern := FirstOf(result.Results(), "Pattern")
	return Cls(pattern.(string)), nil
}

func (p pegHandler) Ref(result *ParseResult) (any, error) {
	_, name := FirstOf(result.Results(), "Name")
	return Ref(name.(string)), nil
}
