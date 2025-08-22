package parser

import (
	"fmt"
	"testing"
)

func TestPegString(t *testing.T) {
	t.Run("Peg String", func(t *testing.T) {
		fmt.Println(Peg())
	})
}

func TestPegParsed(t *testing.T) {
	t.Run("Peg Equal Bootstrap", func(t *testing.T) {
		for rule := range Peg().grammar.Rules() {
			AssertEqual(t, rule.String(), Bootstrap.grammar.Rule(rule.name).String())
		}
	})
}
