package happy_test

import (
	"reflect"
	"testing"

	"github.com/fuwjax/gopase/happy"
)

func TestCompileKeyName(t *testing.T) {
	root := "KeyName"
	tests := []struct {
		name     string
		input    string
		expected any
		err      string
	}{
		{"String key", `"name"`, happy.Lit("name"), ""},
		{"Plain key", `name`, happy.Bracket("name"), ""},
		{"Dotted key", `person.name`, happy.Dotted(happy.Bracket("person"), happy.Bracket("name")), ""},
		{"Bracket key", `person[name]`, happy.Bracket("person", happy.Bracket("name")), ""},
		{"Crazy key", `cities["Europe"]["Amsterdam"].people.0.address[division.name].upper[.]`, happy.Dotted(
			happy.Bracket("cities", happy.Lit("Europe"), happy.Lit("Amsterdam")),
			happy.Bracket("people"),
			happy.Bracket("0"),
			happy.Bracket("address", happy.Dotted(happy.Bracket("division"), happy.Bracket("name"))),
			happy.Bracket("upper", happy.Dot())), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := happy.ParserFrom()(root, tt.input)
			if err == nil && tt.err != "" {
				t.Errorf("err was nil, expected %s", tt.err)
			}
			if err != nil && tt.err == "" {
				t.Errorf("err was %e, expected nil", err)
			}
			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("results was %s, expected %s", results, tt.expected)
			}
		})
	}
}

func TestCompileTags(t *testing.T) {
	root := "Tag"
	tests := []struct {
		name     string
		input    string
		expected any
		err      string
	}{
		{"Section", `(^*name^)(^.^)(^/^)`, &happy.Section{happy.Bracket("name"), &happy.Template{[]happy.Renderer{&happy.Reference{happy.Dot()}}}}, ""},
		{"Comment", `(^# a comment^)`, nil, ""},
		{"Value", `(^name^)`, &happy.Reference{happy.Bracket("name")}, ""},
		{"Value consumes ws", `    ( ^name^ )      `, &happy.Reference{happy.Bracket("name")}, ""},
		{"Else", `(^!name^)`, &happy.Reference{happy.Invert(happy.Bracket("name"))}, ""},
		{"Include", `(^>name^)`, &happy.Include{happy.Bracket("name")}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := happy.ParserFrom()(root, tt.input)
			if err == nil && tt.err != "" {
				t.Errorf("err was nil, expected %s", tt.err)
			}
			if err != nil && tt.err == "" {
				t.Errorf("err was %e, expected nil", err)
			}
			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("results was %s, expected %s", results, tt.expected)
			}
		})
	}
}
