package sample

import (
	"fmt"
	"testing"

	"github.com/fuwjax/gopase/funki/testi"
	"github.com/fuwjax/gopase/parser"
)

func TestPegString(t *testing.T) {
	t.Run("Peg String", func(t *testing.T) {
		fmt.Println(Peg)
	})
}

func TestPegParsed(t *testing.T) {
	t.Run("Peg Equal Bootstrap", func(t *testing.T) {
		testi.AssertEqual(t, Peg.String(), parser.PegGrammar.String())
	})
}
