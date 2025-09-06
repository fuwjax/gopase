package sample_test

import (
	"fmt"
	"testing"

	"github.com/fuwjax/gopase/funki/testi"
	"github.com/fuwjax/gopase/happy/sample"
	"github.com/fuwjax/gopase/parser"
)

func TestPegTemplate(t *testing.T) {
	t.Run("PegTemplate", func(t *testing.T) {
		output, err := sample.RenderPeg(parser.PegGrammar, map[string]any{"package": "sample", "name": "Peg"})
		testi.AssertNil(t, err)
		fmt.Println(output)
	})
}
