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
		result, err := Bootstrap.ParseFrom("EOF", "")
		AssertEqual(t, result, "")
		AssertNil(t, err)

		result, err = Bootstrap.ParseFrom("EOF", "a")
		AssertNil(t, result)
		AssertError(t, err, "at 'a' 1:1 (1) expected not something")
	})
}

func TestBootstrapEol(t *testing.T) {
	t.Run("Bootstrap EOL", func(t *testing.T) {
		result, err := Bootstrap.ParseFrom("EOL", "\n")
		AssertEqual(t, result, "\n")
		AssertNil(t, err)

		result, err = Bootstrap.ParseFrom("EOL", "a")
		AssertNil(t, result)
		AssertError(t, err, "at 'a' 1:1 (1) expected [\\n\\r]")
	})
}
