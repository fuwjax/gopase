package sample_test

import (
	"testing"

	"github.com/fuwjax/gopase/sample"
	"github.com/fuwjax/gopase/when"
)

func TestJson(t *testing.T) {
	// parseJson := input -> () -> sample.ParseJson(input)
	parseJson := func(input string) func() (any, error) {
		return func() (any, error) {
			return sample.ParseJson(input)
		}
	}
	when.YouDoErr("Json String", parseJson(`"abcd"`)).Expect(t, "abcd")
	when.YouDoErr("Json Escape", parseJson(`"\n"`)).Expect(t, "\n")
	when.YouDoErr("Json Number", parseJson(`3.4`)).Expect(t, 3.4)
	when.YouDoErr("Json Array", parseJson(`[1,2,3.4]`)).Expect(t, []any{1.0, 2.0, 3.4})
	when.YouDoErr("Json Object", parseJson(`{"A":"a","B":"b","C":"c"}`)).
		Expect(t, map[string]any{"A": "a", "B": "b", "C": "c"})
	when.YouDoErr("Json Multiline Object", parseJson(`{
			"A": "a",
			"B": "b",
			"C": "c"
		}`)).Expect(t, map[string]any{"A": "a", "B": "b", "C": "c"})
}
