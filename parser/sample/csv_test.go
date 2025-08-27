package sample

import (
	"testing"

	"github.com/fuwjax/gopase/parser"
)

func TestCsvRecords(t *testing.T) {
	t.Run("Csv Records", func(t *testing.T) {
		results, err := ParseCsv(`
A,B,C
a,b,c
`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, [][]string{{"A", "B", "C"}, {"a", "b", "c"}})
	})
}
func TestCsvQuotedRecords(t *testing.T) {
	t.Run("Csv Quoted Records", func(t *testing.T) {
		results, err := ParseCsv(`
A,B,C
"a","b","c"
`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, [][]string{{"A", "B", "C"}, {"a", "b", "c"}})
	})
}
func TestCsvMapRecords(t *testing.T) {
	t.Run("Csv Map Records", func(t *testing.T) {
		results, err := ParseCsvMap(`
A,B,C
a,b,c
`)
		parser.AssertNil(t, err)
		parser.AssertEqual(t, results, []map[string]string{{"A": "a", "B": "b", "C": "c"}})
	})
}
