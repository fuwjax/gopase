package sample

import (
	"iter"
	"strings"
	"sync"

	"github.com/fuwjax/gopase/funki"
	"github.com/fuwjax/gopase/parser"
)

const csvGrammar = `
Records = EOL* Record (EOL+ Record)* EOL* EOF
Record = !EOF Field (',' Field)* 
Field = Quoted / Bare
Quoted = WS '"' Inner '"' WS
Inner = ([^"] / '""')*
Bare = [^,\n\r]*
WS = [ \t]*
EOL = [\n\r]
EOF = !.
`

var CsvParser = sync.OnceValue(func() parser.Parser[[][]string] {
	return parser.NewParser[[][]string]("Records", csvGrammar, csvHandler{})
})

func ParseCsv(input string) ([][]string, error) {
	return CsvParser()(input)
}

func ParseCsvMap(input string) ([]map[string]string, error) {
	records, err := ParseCsv(input)
	if err != nil {
		return nil, err
	}
	header := records[0]
	results := make([]map[string]string, len(records)-1)
	for r, record := range records[1:] {
		results[r] = make(map[string]string)
		for i, key := range header {
			results[r][key] = record[i]
		}
	}
	return results, nil
}

type csvHandler struct{}

func (h csvHandler) Records(results iter.Seq2[string, any]) (any, error) {
	records := funki.ListOf[[]string](results, "Record")
	return records, nil
}

func (h csvHandler) Record(results iter.Seq2[string, any]) (any, error) {
	fields := funki.ListOf[string](results, "Field")
	return fields, nil
}

func (h csvHandler) Field(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Quoted", "Bare")
	return value.(string), nil
}

func (h csvHandler) Quoted(results iter.Seq2[string, any]) (any, error) {
	_, value := funki.FirstOf(results, "Inner")
	value = strings.ReplaceAll(value.(string), "\"\"", "\"")
	return value.(string), nil
}
