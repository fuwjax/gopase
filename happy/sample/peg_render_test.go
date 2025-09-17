package sample_test

import (
	"os"
	"slices"
	"testing"

	"github.com/fuwjax/gopase/happy/sample"
	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/when"
)

func TestPegTemplate(t *testing.T) {
	t.Run("PegTemplate", func(t *testing.T) {
		params := map[string]any{"package": "sample", "name": "Peg"}
		contents := when.YouErr(os.ReadFile("peg.go.txt")).ExpectSuccess(t)
		expected := slices.Collect(parser.Graphemes(string(contents)))

		when.YouErr(sample.RenderPeg(parser.PegGrammar(), params)).
			ExpectMatch(t, func(t *testing.T, actual string) {
				graphemes := slices.Collect(parser.Graphemes(actual))
				when.AssertSlices(t, graphemes, expected)
			})
	})
}
