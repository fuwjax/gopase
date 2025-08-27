package sample

import (
	"iter"
	"strings"

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

func ParseCsv(input string) ([][]string, error) {
	grammar, err := parser.Bootstrap(csvGrammar)
	if err != nil {
		return nil, err
	}
	result, err := parser.Parse("Records", grammar, parser.WrapHandler(csvHandler{}), input)
	if err != nil {
		return nil, err
	}
	return result.([][]string), nil
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
	return parser.Cast[[]string](parser.ListOf(results, "Record")), nil
}

func (h csvHandler) Record(results iter.Seq2[string, any]) (any, error) {
	return parser.Cast[string](parser.ListOf(results, "Field")), nil
}

func (h csvHandler) Field(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Quoted", "Bare")
	return value.(string), nil
}

func (h csvHandler) Quoted(results iter.Seq2[string, any]) (any, error) {
	_, value := parser.FirstOf(results, "Inner")
	value = strings.ReplaceAll(value.(string), "\"\"", "\"")
	return value.(string), nil
}
