package parser

import (
	"testing"
)

func TestParserBasic(t *testing.T) {
	t.Run("Parser Basic", func(t *testing.T) {
		parser := NewGrammar().AddRule("S", Lit("abc")).Parser("S", nil)
		result, err := parser.Parse("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserReuse(t *testing.T) {
	t.Run("Parser Reuse", func(t *testing.T) {
		parser := NewGrammar().AddRule("S", Lit("abc")).Parser("S", nil)
		parser.Parse("abc")
		result, err := parser.Parse("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserSequence(t *testing.T) {
	t.Run("Parser Sequence", func(t *testing.T) {
		parser := NewGrammar().AddRule("S", Seq(Lit("a"), Lit("b"), Lit("c"))).Parser("S", nil)
		result, err := parser.Parse("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserRef(t *testing.T) {
	t.Run("Parser Ref", func(t *testing.T) {
		parser := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc")).Parser("S", nil)
		result, err := parser.Parse("abc")
		AssertNil(t, err)
		AssertEqual(t, result, "abc")
	})
}
func TestParserHandler(t *testing.T) {
	t.Run("Parser Hander", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result *ParseResult) (any, error) { return ListOf(result.Results(), "T"), nil }
		parser := NewGrammar().AddRule("S", Ref("T")).AddRule("T", Lit("abc")).Parser("S", handler)
		result, err := parser.Parse("abc")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"abc"})
	})
}
func TestParserHandlerArray(t *testing.T) {
	t.Run("Parser Hander Array", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result *ParseResult) (any, error) { return ListOf(result.Results(), "T"), nil }
		parser := NewGrammar().AddRule("S", Seq(Ref("T"), Ref("T"))).AddRule("T", Lit("ab")).Parser("S", handler)
		result, err := parser.Parse("abab")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})
	})
}
