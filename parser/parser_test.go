package parser

import (
	"iter"
	"testing"
)

func TestParserBasic(t *testing.T) {
	t.Run("Parser Basic", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Lit("abc"))
		parser := NewParser[any]("S", grammar, nil)
		result, err := parser("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserReuse(t *testing.T) {
	t.Run("Parser Reuse", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Lit("abc"))
		parser := NewParser[any]("S", grammar, nil)
		parser("abc")
		result, err := parser("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserSequence(t *testing.T) {
	t.Run("Parser Sequence", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Seq(Lit("a"), Lit("b"), Lit("c")))
		parser := NewParser[any]("S", grammar, nil)
		result, err := parser("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserRef(t *testing.T) {
	t.Run("Parser Ref", func(t *testing.T) {
		grammar := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc"))
		parser := NewParser[any]("S", grammar, nil)
		result, err := parser("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserHandler(t *testing.T) {
	t.Run("Parser Hander", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc"))
		parser := NewParser[any]("S", grammar, handler)
		result, err := parser("abc")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"abc"})
	})
}
func TestParserHandlerArray(t *testing.T) {
	t.Run("Parser Hander Array", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar().AddRule("S", Seq(Ref("T"), Ref("T"))).AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)
		result, err := parser("abab")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})
	})
}
