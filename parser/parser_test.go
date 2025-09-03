package parser

import (
	"iter"
	"slices"
	"testing"

	"github.com/fuwjax/gopase/funki"
	"github.com/fuwjax/gopase/funki/testi"
)

func TestParserBasic(t *testing.T) {
	t.Run("Parser Basic", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Lit("abc"))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(nil))
		result, err := parser("abc")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, "abc")
	})
}
func TestParserReuse(t *testing.T) {
	t.Run("Parser Reuse", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Lit("abc"))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(nil))
		parser("abc")
		result, err := parser("abc")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, "abc")
	})
}
func TestParserSequence(t *testing.T) {
	t.Run("Parser Sequence", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Seq(Lit("a"), Lit("b"), Lit("c")))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(nil))
		result, err := parser("abc")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, "abc")
	})
}
func TestParserRef(t *testing.T) {
	t.Run("Parser Ref", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc"))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(nil))
		result, err := parser("abc")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, "abc")
	})
}
func TestParserHandler(t *testing.T) {
	t.Run("Parser Hander", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
			return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
		}
		grammar := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc"))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(handler))
		result, err := parser("abc")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, []any{"abc"})
	})
}
func TestParserHandlerArray(t *testing.T) {
	t.Run("Parser Hander Array", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) {
			return slices.Collect(funki.Values(funki.FilterKeys(result, "T"))), nil
		}
		grammar := NewGrammar().AddRule("S", Seq(Ref("T"), Ref("T"))).AddRule("T", Lit("ab"))
		parser := BootstrapParser[any]("S", grammar, WrapHandler(handler))
		result, err := parser("abab")
		testi.AssertNil(t, err)
		testi.AssertEqual(t, result, []any{"ab", "ab"})
	})
}
