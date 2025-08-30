package sample

import (
	"fmt"
	"testing"

	"github.com/fuwjax/gopase/parser"
)

func TestJsonGrammar(t *testing.T) {
	t.Run("Json Grammar String", func(t *testing.T) {
		grammar, err := JsonGrammar()
		parser.AssertNil(t, err)
		fmt.Println(grammar)
	})
}
func TestJsonString(t *testing.T) {
	t.Run("Json String", func(t *testing.T) {
		results, err := ParseJson(`"abcd"`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, "abcd")
	})
	t.Run("Json Escape", func(t *testing.T) {
		results, err := ParseJsonFrom("String", `"\n"`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, "\n")
	})
}

func TestJsonNumber(t *testing.T) {
	t.Run("Json Number", func(t *testing.T) {
		results, err := ParseJson(`3.4`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, 3.4)
	})
}

func TestJsonArray(t *testing.T) {
	t.Run("Json Array", func(t *testing.T) {
		results, err := ParseJson(`[1,2,3.4]`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, []any{1.0, 2.0, 3.4})
	})
}

func TestJsonObject(t *testing.T) {
	t.Run("Json Object", func(t *testing.T) {
		results, err := ParseJson(`{"A":"a","B":"b","C":"c"}`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, map[string]any{"A": "a", "B": "b", "C": "c"})
	})
}

func TestJsonMultilineObject(t *testing.T) {
	t.Run("Json Multiline Object", func(t *testing.T) {
		results, err := ParseJson(`{
			"A": "a",
			"B": "b",
			"C": "c"
		}`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, map[string]any{"A": "a", "B": "b", "C": "c"})
	})
}
