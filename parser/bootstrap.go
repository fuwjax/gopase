package parser

import (
	"iter"
	"slices"
	"strings"

	"github.com/fuwjax/gopase/funki"
)

/*
The Peg-grammar parser.
*/
var Bootstrap = BootstrapParser[*Grammar]("Grammar", PegGrammar(), PegHandler)

var BootstrapFrom = BootstrapParserFrom(PegGrammar(), PegHandler)

var PegHandler = WrapHandler(pegHandler{})

type pegHandler struct{}

func (p pegHandler) Grammar(result iter.Seq2[string, any]) (any, error) {
	rules := slices.Collect(funki.Cast[*Rule](funki.FilterNonNil(funki.Values(funki.FilterKeys(result, "Line")))))
	grammar := NewGrammar()
	for _, rule := range rules {
		grammar.Add(rule)
	}
	return grammar, nil
}

func (p pegHandler) Line(result iter.Seq2[string, any]) (any, error) {
	_, rule := funki.FirstOf(result, "Rule")
	return rule, nil
}

func (p pegHandler) Rule(result iter.Seq2[string, any]) (any, error) {
	_, name := funki.FirstOf(result, "Name")
	_, expr := funki.FirstOf(result, "Expr")
	return NewRule(name.(string), expr.(Expr)), nil
}

func (p pegHandler) Expr(result iter.Seq2[string, any]) (any, error) {
	seqs := funki.ListOf[Expr](result, "Seq")
	return Alt(seqs...), nil
}

func (p pegHandler) Seq(result iter.Seq2[string, any]) (any, error) {
	prefixs := funki.ListOf[Expr](result, "Prefix")
	return Seq(prefixs...), nil
}

func (p pegHandler) Prefix(result iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(result, "AndExpr", "NotExpr", "Suffix")
	return value, nil
}

func (p pegHandler) AndExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Suffix")
	return See(expr.(Expr)), nil
}

func (p pegHandler) NotExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Suffix")
	return Not(expr.(Expr)), nil
}

func (p pegHandler) Suffix(result iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(result, "OptExpr", "RepExpr", "ReqExpr", "Primary")
	return value, nil
}

func (p pegHandler) OptExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Primary")
	return Opt(expr.(Expr)), nil
}

func (p pegHandler) RepExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Primary")
	return Rep(expr.(Expr)), nil
}

func (p pegHandler) ReqExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Primary")
	return Req(expr.(Expr)), nil
}

func (p pegHandler) Primary(result iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(result, "Dot", "ParExpr", "Literal", "CharClass", "Ref")
	return value, nil
}

func (p pegHandler) Dot(result iter.Seq2[string, any]) (any, error) {
	return Dot(), nil
}

func (p pegHandler) ParExpr(result iter.Seq2[string, any]) (any, error) {
	_, expr := funki.FirstOf(result, "Expr")
	return expr, nil
}

func (p pegHandler) Literal(result iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(result, "SingleLit", "DoubleLit")
	return Lit(value.(string)), nil
}

func (p pegHandler) SingleLit(results iter.Seq2[string, any]) (any, error) {
	var sb strings.Builder
	for name, result := range results {
		switch name {
		case "SinglePlain":
			sb.WriteString(result.(string))
		case "SingleEscape":
			switch result.(string) {
			case "n":
				sb.WriteString("\n")
			case "r":
				sb.WriteString("\r")
			case "t":
				sb.WriteString("\t")
			default:
				sb.WriteString(result.(string))
			}
		}
	}
	return sb.String(), nil
}

func (p pegHandler) DoubleLit(results iter.Seq2[string, any]) (any, error) {
	var sb strings.Builder
	for name, result := range results {
		switch name {
		case "DoublePlain":
			sb.WriteString(result.(string))
		case "DoubleEscape":
			switch result.(string) {
			case "n":
				sb.WriteString("\n")
			case "r":
				sb.WriteString("\r")
			case "t":
				sb.WriteString("\t")
			default:
				sb.WriteString(result.(string))
			}
		}
	}
	return sb.String(), nil
}

func (p pegHandler) CharClass(result iter.Seq2[string, any]) (any, error) {
	_, pattern := funki.FirstOf(result, "Pattern")
	return Cls(pattern.(string)), nil
}

func (p pegHandler) Ref(result iter.Seq2[string, any]) (any, error) {
	_, name := funki.FirstOf(result, "Name")
	return Ref(name.(string)), nil
}
