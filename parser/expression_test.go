package parser_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/fuwjax/gopase/funki"
	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/when"
)

func testParser(parser parser.Parser[any], input string) when.WhenOpErr[any] {
	return func() (any, error) {
		return parser(input)
	}
}

func TestParserSeq(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.Ref("T"), parser.Ref("T")))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("seq matching", testParser(parser, "abab")).Expect(t, []any{"ab", "ab"})
	when.YouDoErr("seq miss", testParser(parser, "abba")).ExpectError(t, "at 'b' 1:3 (3) expected a\nwhile in T\nwhile in S")
}

func TestParserAlt(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Alt(parser.Seq(parser.Ref("T"), parser.Lit("a")), parser.Seq(parser.Ref("T"), parser.Lit("b"))))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("alt matching repeat", testParser(parser, "abab")).Expect(t, []any{"ab"})
	when.YouDoErr("alt matching", testParser(parser, "abba")).Expect(t, []any{"ab"})
	when.YouDoErr("alt miss", testParser(parser, "acba")).ExpectError(t, "at 'c' 1:2 (2) expected b\nwhile in T\nat 'c' 1:2 (2) expected b\nwhile in T\nwhile in S")
}

func TestParserCls(t *testing.T) {
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Cls("[a-f]"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))

	when.YouDoErr("cls matching", testParser(parser, "a")).Expect(t, "a")
	when.YouDoErr("cls miss", testParser(parser, "x")).ExpectError(t, "at 'x' 1:1 (1) expected [a-f]\nwhile in S")
}

func TestParserDot(t *testing.T) {
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Dot())
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))

	when.YouDoErr("dot matching", testParser(parser, "a")).Expect(t, "a")
	when.YouDoErr("dot eof", testParser(parser, "")).ExpectError(t, "at EOF 1:0 (0) expected anything\nwhile in S")
}

func TestParserOpt(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.Opt(parser.Ref("T")), parser.Lit("a")))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("opt matching", testParser(parser, "aba")).Expect(t, []any{"ab"})
	when.YouDoErr("opt miss", testParser(parser, "a")).Expect(t, []any(nil))
}

func TestParserRep(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.Rep(parser.Ref("T")), parser.Lit("a")))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("rep matching", testParser(parser, "aba")).Expect(t, []any{"ab"})
	when.YouDoErr("rep miss", testParser(parser, "abab")).ExpectError(t, "at EOF 1:5 (5) expected a\nwhile in S")
	when.YouDoErr("rep matching multiple", testParser(parser, "ababa")).Expect(t, []any{"ab", "ab"})
	when.YouDoErr("rep empty", testParser(parser, "a")).Expect(t, []any(nil))
}

func TestParserReq(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.Req(parser.Ref("T")), parser.Lit("a")))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("req matching", testParser(parser, "aba")).Expect(t, []any{"ab"})
	when.YouDoErr("req miss", testParser(parser, "abab")).ExpectError(t, "at EOF 1:5 (5) expected a\nwhile in S")
	when.YouDoErr("req matching multiple", testParser(parser, "ababa")).Expect(t, []any{"ab", "ab"})
	when.YouDoErr("req empty", testParser(parser, "a")).ExpectError(t, "at EOF 1:2 (2) expected b\nwhile in T\nwhile in S")
}

func TestParserSee(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.See(parser.Ref("T")), parser.Lit("a")))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("see matching", testParser(parser, "aba")).Expect(t, []any(nil))
	when.YouDoErr("see miss", testParser(parser, "bab")).ExpectError(t, "at 'b' 1:1 (1) expected a\nwhile in T\nwhile in S")
}

func TestParserNot(t *testing.T) {
	handler := make(map[string]parser.Converter)
	handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
		return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
	}
	grammar := parser.NewGrammar()
	grammar.AddRule("S", parser.Seq(parser.Req(parser.Ref("T")), parser.Lit("a"), parser.Not(parser.Dot())))
	grammar.AddRule("T", parser.Lit("ab"))
	parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))

	when.YouDoErr("not matching", testParser(parser, "ababa")).Expect(t, []any{"ab", "ab"})
	when.YouDoErr("not miss", testParser(parser, "ababaa")).ExpectError(t, "at 'a' 1:6 (6) expected not something\nwhile in S")
}
