package sample_test

import (
	"os"
	"slices"
	"testing"

	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/sample"
	"github.com/fuwjax/gopase/when"
)

func AssertToken(t *testing.T, actual, expected *parser.Grapheme) bool {
	if actual.Token != expected.Token {
		return when.AssertEqual(t, actual, expected)
	}
	return true
}

func MatchGraphemes(expected string) when.Matcher[string] {
	expectedGs := slices.Collect(parser.Graphemes(expected))
	return func(t *testing.T, actual string) bool {
		graphemes := parser.Graphemes(actual)
		return when.MatchSeq(AssertToken, expectedGs...)(t, graphemes)
	}
}

func TestPegTemplate(t *testing.T) {
	t.Run("PegTemplate", func(t *testing.T) {
		params := map[string]any{"package": "parser", "name": "Peg", "inPackage": true}
		contents := when.YouErr(os.ReadFile("peg.gold")).ExpectSuccess(t)

		when.YouErr(sample.RenderPeg(parser.PegGrammar(), params)).
			//Expect(t, string(contents))
			ExpectMatch(t, MatchGraphemes(string(contents)))
	})
}
