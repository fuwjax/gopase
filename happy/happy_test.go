package happy_test

import (
	"testing"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/parser/sample"
)

func TestHappyInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     any
		expected string
		err      string
	}{
		{"Plaintext", "This is plain text!", nil, "This is plain text!", ""},
		{"Kitty", "This is (^.^) text!", "kitty", "This is kitty text!", ""},
		{"Interpolate", "This is (^name^) text!", map[string]any{"name": "Bob"}, "This is Bob text!", ""},
		{"Strip leading space", `This is 
		
		( ^name^) text!`, map[string]any{"name": "sticky"}, "This issticky text!", ""},
		{"Strip trailing space", `This is    (^name^ ) text!`, map[string]any{"name": "plain"}, "This is    plaintext!", ""},
		{"Handle inner whitespace", `This is (^
		name
		^) text!`, map[string]any{"name": "valid"}, "This is valid text!", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := happy.Render(tt.template, tt.data, nil)
			if err == nil && tt.err != "" {
				t.Errorf("err was nil, expected %s", tt.err)
			}
			if err != nil && tt.err == "" {
				t.Errorf("err was %e, expected nil", err)
			}
			if results != tt.expected {
				t.Errorf("results was %s, expected %s", results, tt.expected)
			}
		})
	}
}

func TestHappyJsonInterpolation(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     string
		expected string
		err      string
	}{
		{"Simple", `
		(^name^)
		(^age^)
		`, `
		{"name":"Bob","age":123}
		`, `
		Bob
		123
		`, ""},
		{"Dotted", `
		(^person.name^)
		(^person.age^)
		`, `
		{"person":{"name":"Bob Hope","age":55}}
		`, `
		Bob Hope
		55
		`, ""},
		{"Bracketted", `
		( ^ person[name] ^ )
		( ^ person[age] ^ )
		`, `
		{"person":{"name":"Bob Ross","age":35}}
		`, `Bob Ross35`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := sample.ParseJson(tt.data)
			if err != nil {
				t.Errorf("json parse err was %e, expected nil", err)
			}
			results, err := happy.Render(tt.template, data, nil)
			if err == nil && tt.err != "" {
				t.Errorf("err was nil, expected %s", tt.err)
			}
			if err != nil && tt.err == "" {
				t.Errorf("err was %e, expected nil", err)
			}
			if results != tt.expected {
				t.Errorf("results was %s, expected %s", results, tt.expected)
			}
		})
	}
}
