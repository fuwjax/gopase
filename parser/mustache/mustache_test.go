package mustache

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/fuwjax/gopase/parser"
	"github.com/fuwjax/gopase/parser/sample"
)

func TestMustacheString(t *testing.T) {
	t.Run("Mustache Grammar String", func(t *testing.T) {
		grammar, err := Grammar()
		parser.AssertNil(t, err)
		fmt.Println(grammar)
	})
}

type TestSuite struct {
	Overview string `jsonName:"overview"`
	Tests    []Test `jsonName:"tests"`
}

type Test struct {
	Name     string `jsonName:"name"`
	Desc     string `jsonName:"desc"`
	Data     any    `jsonName:"data"`
	Template string `jsonName:"template"`
	Expected string `jsonName:"expected"`
}

func mustacheTest(t *testing.T, filename string) {
	resp, err := http.Get(baseUrl + filename)
	parser.AssertNil(t, err)
	defer resp.Body.Close()
	parser.AssertEqual(t, resp.StatusCode, http.StatusOK)
	bytes, err := io.ReadAll(resp.Body)
	parser.AssertNil(t, err)
	results, err := sample.ParseJson(string(bytes))
	parser.AssertNil(t, err)
	suite, err := sample.ConvertJson[TestSuite](results)
	parser.AssertNil(t, err)
	fmt.Println(suite.Overview)
	for _, test := range suite.Tests {
		t.Run(test.Name, func(t *testing.T) {
			fmt.Println(test.Desc)
			result, err := Render(test.Template, test.Data)
			parser.AssertNil(t, err)
			parser.AssertEqual(t, result, test.Expected)
		})
	}
}

const baseUrl = "https://raw.githubusercontent.com/mustache/spec/refs/heads/master/specs/"

/*
comments.json
interpolation.json
sections.json
partials.json
inverted.json
delimiters.json
~dynamic-names.json
~inheritance.json
*/

func TestComments(t *testing.T) {
	mustacheTest(t, "comments.json")
}

func TestInterpolation(t *testing.T) {
	mustacheTest(t, "interpolation.json")
}
