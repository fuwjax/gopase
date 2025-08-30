package parser

import (
	"iter"
	"testing"
)

func TestParserSeq(t *testing.T) {
	t.Run("Parser Seq", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(Ref("T"), Ref("T")))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("abab")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})

		result, err = parser("abba")
		AssertError(t, err, "at 'b' 1:3 (3) expected a\nwhile in T\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserAlt(t *testing.T) {
	t.Run("Parser Alt", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Alt(Seq(Ref("T"), Lit("a")), Seq(Ref("T"), Lit("b"))))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("abab")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab"})

		result, err = parser("abba")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab"})

		result, err = parser("acba")
		AssertError(t, err, "at 'c' 1:2 (2) expected b\nwhile in T\nat 'c' 1:2 (2) expected b\nwhile in T\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserCls(t *testing.T) {
	t.Run("Parser Cls", func(t *testing.T) {
		grammar := NewGrammar()
		grammar.AddRule("S", Cls("[a-f]"))
		parser := NewParser[any]("S", grammar, nil)

		result, err := parser("a")
		AssertNil(t, err)
		AssertEqual(t, result, "a")

		result, err = parser("x")
		AssertError(t, err, "at 'x' 1:1 (1) expected [a-f]\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserDot(t *testing.T) {
	t.Run("Parser Dot", func(t *testing.T) {
		grammar := NewGrammar()
		grammar.AddRule("S", Dot())
		parser := NewParser[any]("S", grammar, nil)

		result, err := parser("a")
		AssertNil(t, err)
		AssertEqual(t, result, "a")

		result, err = parser("")
		AssertError(t, err, "at EOF 1:0 (0) expected anything\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserOpt(t *testing.T) {
	t.Run("Parser Opt", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(Opt(Ref("T")), Lit("a")))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("aba")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab"})

		result, err = parser("a")
		AssertNil(t, err)
		AssertEqual(t, result, []any{})
	})
}

func TestParserRep(t *testing.T) {
	t.Run("Parser Rep", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(Rep(Ref("T")), Lit("a")))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("aba")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab"})

		result, err = parser("abab")
		AssertError(t, err, "at EOF 1:5 (5) expected a\nwhile in S")
		AssertNil(t, result)

		result, err = parser("ababa")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})

		result, err = parser("a")
		AssertNil(t, err)
		AssertEqual(t, result, []any{})
	})
}

func TestParserReq(t *testing.T) {
	t.Run("Parser Req", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(Req(Ref("T")), Lit("a")))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("aba")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab"})

		result, err = parser("abab")
		AssertError(t, err, "at EOF 1:5 (5) expected a\nwhile in S")
		AssertNil(t, result)

		result, err = parser("ababa")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})

		result, err = parser("a")
		AssertError(t, err, "at EOF 1:2 (2) expected b\nwhile in T\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserSee(t *testing.T) {
	t.Run("Parser See", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(See(Ref("T")), Lit("a")))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("aba")
		AssertNil(t, err)
		AssertEqual(t, result, []any{})

		result, err = parser("bab")
		AssertError(t, err, "at 'b' 1:1 (1) expected a\nwhile in T\nwhile in S")
		AssertNil(t, result)
	})
}

func TestParserNot(t *testing.T) {
	t.Run("Parser Not", func(t *testing.T) {
		handler := make(map[string]Converter)
		handler["S"] = func(result iter.Seq2[string, any]) (any, error) { return ListOf(result, "T"), nil }
		grammar := NewGrammar()
		grammar.AddRule("S", Seq(Req(Ref("T")), Lit("a"), Not(Dot())))
		grammar.AddRule("T", Lit("ab"))
		parser := NewParser[any]("S", grammar, handler)

		result, err := parser("ababaa")
		AssertError(t, err, "at 'a' 1:6 (6) expected not something\nwhile in S")
		AssertNil(t, result)

		result, err = parser("ababa")
		AssertNil(t, err)
		AssertEqual(t, result, []any{"ab", "ab"})
	})
}
