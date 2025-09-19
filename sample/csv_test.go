package sample_test

import (
	"testing"

	"github.com/fuwjax/gopase/sample"
	"github.com/fuwjax/gopase/when"
)

func TestCsvRecords(t *testing.T) {
	t.Run("Csv Records", func(t *testing.T) {
		when.YouErr(sample.ParseCsv(`
A,B,C
a,b,c
`)).Expect(t, [][]string{{"A", "B", "C"}, {"a", "b", "c"}})
	})
}
func TestCsvQuotedRecords(t *testing.T) {
	t.Run("Csv Quoted Records", func(t *testing.T) {
		when.YouErr(sample.ParseCsv(`
A,B,C
"a","b","c"
`)).Expect(t, [][]string{{"A", "B", "C"}, {"a", "b", "c"}})
	})
}
func TestCsvMapRecords(t *testing.T) {
	t.Run("Csv Map Records", func(t *testing.T) {
		when.YouErr(sample.ParseCsvMap(`
A,B,C
a,b,c
`)).Expect(t, []map[string]string{{"A": "a", "B": "b", "C": "c"}})
	})
}
