package parser_test

import (
	"iter"
	"slices"
	"testing"

	"github.com/fuwjax/gopase/funki"
	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/when"
)

func TestParserBasic(t *testing.T) {
	t.Run("Parser Basic", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Lit("abc"))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("abc")).Expect(t, "abc")
	})
}
func TestParserReuse(t *testing.T) {
	t.Run("Parser Reuse", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Lit("abc"))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		parser("abc")
		when.YouErr(parser("abc")).Expect(t, "abc")
	})
}
func TestParserSequence(t *testing.T) {
	t.Run("Parser Sequence", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Seq(parser.Lit("a"), parser.Lit("b"), parser.Lit("c")))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("abc")).Expect(t, "abc")
	})
}
func TestParserRef(t *testing.T) {
	t.Run("Parser Ref", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Ref("T")).AddRule("T", parser.Lit("abc"))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("abc")).Expect(t, "abc")
	})
}
func TestParserHandler(t *testing.T) {
	t.Run("Parser Hander", func(t *testing.T) {
		handler := make(map[string]parser.Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
			return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
		}
		grammar := parser.NewGrammar().AddRule("S", parser.Ref("T")).AddRule("T", parser.Lit("abc"))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))
		when.YouErr(parser("abc")).Expect(t, []any{"abc"})
	})
}
func TestParserHandlerArray(t *testing.T) {
	t.Run("Parser Hander Array", func(t *testing.T) {
		handler := make(map[string]parser.Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
			return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
		}
		grammar := parser.NewGrammar().AddRule("S", parser.Seq(parser.Ref("T"), parser.Ref("T"))).AddRule("T", parser.Lit("ab"))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(handler))
		when.YouErr(parser("abab")).Expect(t, []any{"ab", "ab"})
	})
}
func TestParserRightRecursion(t *testing.T) {
	t.Run("Parser Right Recursion", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Alt(parser.Seq(parser.Lit("a"), parser.Ref("S")), parser.Lit("a")))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("a")).Expect(t, "a")
		when.YouErr(parser("aaa")).Expect(t, "aaa")
		when.YouErr(parser("aab")).Expect(t, "aa")
	})
}

func TestParserLeftRecursion(t *testing.T) {
	t.Run("Parser Left Recursion", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Alt(parser.Seq(parser.Ref("S"), parser.Lit("a")), parser.Lit("a")))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("a")).Expect(t, "a")
		when.YouErr(parser("aaa")).Expect(t, "aaa")
		when.YouErr(parser("aab")).Expect(t, "aa")
	})
}

func TestParserIndirectLeftRecursion(t *testing.T) {
	t.Run("Parser Left Recursion", func(t *testing.T) {
		grammar := parser.NewGrammar().AddRule("S", parser.Alt(parser.Seq(parser.Ref("T"), parser.Lit("b")), parser.Lit("a")))
		grammar.AddRule("T", parser.Alt(parser.Seq(parser.Ref("S"), parser.Lit("a")), parser.Lit("c")))
		parser := parser.BootstrapParser[any]("S", grammar, parser.WrapHandler(nil))
		when.YouErr(parser("a")).Expect(t, "a")
		when.YouErr(parser("aaba")).Expect(t, "aab")
		when.YouErr(parser("aabab")).Expect(t, "aabab")
	})
}
