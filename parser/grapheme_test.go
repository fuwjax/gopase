package parser_test

import (
	"iter"
	"testing"

	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/when"
)

func TestNew(t *testing.T) {
	newG := func(input string) when.WhenOp[*parser.Grapheme] {
		return func() *parser.Grapheme {
			return parser.NewGrapheme(input)
		}
	}
	when.YouDo("Normal New", newG("abc")).Expect(t, parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))
	when.YouDo("Empty New", newG("")).Expect(t, parser.NewTestGrapheme("", "", 1, 0, 0, -1, 0))
	when.YouDo("New Line New", newG("\nabc")).Expect(t, parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))
}

func TestGraphemeNext(t *testing.T) {
	next := func(g *parser.Grapheme) when.WhenOp[*parser.Grapheme] {
		return func() *parser.Grapheme {
			return g.Next()
		}
	}
	when.YouDo("Initial Next", next(parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))).
		Expect(t, parser.NewTestGrapheme("b", "c", 1, 2, 2, 2320992, 16))
	when.YouDo("Normal Next", next(parser.NewTestGrapheme("a", "bc", 3, 17, 41, 2320992, 16))).
		Expect(t, parser.NewTestGrapheme("b", "c", 3, 18, 42, 2320992, 16))
	when.YouDo("End Next", next(parser.NewTestGrapheme("a", "", 7, 1, 53, 2320992, 16))).
		Expect(t, parser.NewTestGrapheme("", "", 7, 2, 54, -1, 0))
	when.YouDo("After End Next", next(parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))).
		Expect(t, parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))
	when.YouDo("New Line Next", next(parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))).
		Expect(t, parser.NewTestGrapheme("a", "bc", 2, 1, 2, 2320992, 16))
}

func TestGraphemeIsEof(t *testing.T) {
	eof := func(g *parser.Grapheme) when.WhenOp[bool] {
		return func() bool {
			return g.IsEof()
		}
	}
	when.YouDo("Initial IsEof", eof(parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("Normal IsEof", eof(parser.NewTestGrapheme("a", "bc", 3, 17, 41, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("End IsEof", eof(parser.NewTestGrapheme("a", "", 7, 1, 53, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("After End IsEof", eof(parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))).
		ExpectSuccess(t)
	when.YouDo("New Line IsEof", eof(parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))).
		ExpectFailure(t)
}

func TestGraphemeIsEol(t *testing.T) {
	eol := func(g *parser.Grapheme) when.WhenOp[bool] {
		return func() bool {
			return g.IsEol()
		}
	}
	when.YouDo("Initial IsEol", eol(parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("Normal IsEol", eol(parser.NewTestGrapheme("a", "bc", 3, 17, 41, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("End IsEol", eol(parser.NewTestGrapheme("a", "", 7, 1, 53, 2320992, 16))).
		ExpectFailure(t)
	when.YouDo("After End IsEol", eol(parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))).
		ExpectFailure(t)
	when.YouDo("New Line IsEol", eol(parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))).
		ExpectSuccess(t)
}

func TestGrapheme_String(t *testing.T) {
	str := func(g *parser.Grapheme) when.WhenOp[string] {
		return func() string {
			return g.String()
		}
	}
	when.YouDo("Initial String", str(parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))).
		Expect(t, "'a' 1:1 (1)")
	when.YouDo("Normal String", str(parser.NewTestGrapheme("a", "bc", 3, 17, 41, 2320992, 16))).
		Expect(t, "'a' 3:17 (41)")
	when.YouDo("End String", str(parser.NewTestGrapheme("a", "", 7, 1, 53, 2320992, 16))).
		Expect(t, "'a' 7:1 (53)")
	when.YouDo("After End String", str(parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))).
		Expect(t, "EOF 5:8 (14)")
	when.YouDo("New Line String", str(parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))).
		Expect(t, "'\n' 1:1 (1)")
}

func TestGrapheme_Error(t *testing.T) {
	err := func(g *parser.Grapheme, expected string) when.WhenOp[string] {
		return func() string {
			return g.Error(expected).Error()
		}
	}
	when.YouDo("Normal Error", err(parser.NewTestGrapheme("a", "bc", 3, 2, 15, 2320992, 16), "b")).
		Expect(t, "at 'a' 3:2 (15) expected b")
}

func TestGraphemes(t *testing.T) {
	graphemes := func(input string) when.WhenOp[iter.Seq[*parser.Grapheme]] {
		return func() iter.Seq[*parser.Grapheme] {
			return parser.Graphemes(input)
		}
	}
	when.YouDo("Normal Graphemes", graphemes("a\nb\\'c")).ExpectMatch(t, when.MatchSeq(
		parser.NewTestGrapheme("a", "\nb\\'c", 1, 1, 1, 8414242, 20),
		parser.NewTestGrapheme("\n", "b\\'c", 1, 2, 2, 2320992, 14),
		parser.NewTestGrapheme("b", "\\'c", 2, 1, 3, 2334720, 20),
		parser.NewTestGrapheme("\\", "'c", 2, 2, 4, 2236416, 20),
		parser.NewTestGrapheme("'", "c", 2, 3, 5, 2320992, 20),
		parser.NewTestGrapheme("c", "", 2, 4, 6, 0, 30)))
}
