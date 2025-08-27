package parser

import (
	"fmt"
	"testing"
)

func TestBootstrapString(t *testing.T) {
	t.Run("Bootstrap String", func(t *testing.T) {
		fmt.Println(Bootstrap)
	})
}

func TestBootstrapEof(t *testing.T) {
	t.Run("Bootstrap EOF", func(t *testing.T) {
		result, err := Parse("EOF", PegGrammar, PegHandler, "")
		AssertEqual(t, result, "")
		AssertNil(t, err)

		result, err = Parse("EOF", PegGrammar, PegHandler, "a")
		AssertNil(t, result)
		AssertError(t, err, "at 'a' 1:1 (1) expected not something")
	})
}

func TestBootstrapEol(t *testing.T) {
	t.Run("Bootstrap EOL", func(t *testing.T) {
		result, err := Parse("EOL", PegGrammar, PegHandler, "\n")
		AssertEqual(t, result, "\n")
		AssertNil(t, err)

		result, err = Parse("EOL", PegGrammar, PegHandler, "a")
		AssertNil(t, result)
		AssertError(t, err, "at 'a' 1:1 (1) expected [\\n\\r]")
	})
}
