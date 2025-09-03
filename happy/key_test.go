package happy_test

import (
	"reflect"
	"testing"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/parser/sample"
)

func checkError(t *testing.T, err error, format string, arg any) {
	if err != nil {
		t.Errorf(format, err, arg)
	}
}

func TestKeyResolve(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		context  string
		expected string
		ok       bool
	}{
		{"String key", `"name"`, `[]`, `"name"`, true},
		{"String key ignores context", `"name"`, `[{"name":"Bob"}]`, `"name"`, true},
		{"Plain key uses context", `name`, `[{"name":"Bob"}]`, `"Bob"`, true},
		{"Plain key fails without context", `name`, `[{"first_name":"Bob"}]`, `null`, false},
		{"Plain key uses latest context", `name`, `[{"name":"Bob"},{"name":"Jim"}]`, `"Jim"`, true},
		{"Plain key travels context", `name`, `[{"name":"Bob"},{"first_name":"Jim"}]`, `"Bob"`, true},
		{"Dot returns context", `.`, `[{"name":"Bob"}]`, `{"name":"Bob"}`, true},
		{"Dotted uses nested context", `person.name`, `[{"person":{"name":"Bob"}}]`, `"Bob"`, true},
		{"Brackets uses nested context", `person["name"]`, `[{"person":{"name":"Bob"}}]`, `"Bob"`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := happy.ParserFrom()("KeyName", tt.key)
			checkError(t, err, "could not parse %s into key", tt.key)
			contextJson, err := sample.ParseJson(tt.context)
			checkError(t, err, "could not parse %s into json", tt.context)
			expected, err := sample.ParseJson(tt.expected)
			checkError(t, err, "could not parse %s into json", tt.expected)

			context := happy.ContextOf(contextJson.([]any)...)
			actual, ok := key.(happy.Key).Resolve(context)
			if ok != tt.ok {
				t.Errorf("ok was %t, expected %t", ok, tt.ok)
			}
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("actual was %s, expected %s", actual, expected)
			}
		})
	}
}

func TestKeyResolveName(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		context  string
		expected string
	}{
		{"String key", `"name"`, `[]`, "name"},
		{"String key ignores context", `"name"`, `[{"name":"Bob"}]`, "name"},
		{"Plain key uses context", `name`, `[{"name":"Bob"}]`, "Bob"},
		{"Plain key reverts to literal without context", `name`, `[{"first_name":"Bob"}]`, "name"},
		{"Plain key uses latest context", `name`, `[{"name":"Bob"},{"name":"Jim"}]`, "Jim"},
		{"Plain key travels context", `name`, `[{"name":"Bob"},{"first_name":"Jim"}]`, "Bob"},
		{"Dot returns String() of current context", `.`, `[{"name":"Bob"}]`, "map[name:Bob]"},
		{"Dotted uses nested context", `person.name`, `[{"person":{"name":"Bob"}}]`, "Bob"},
		{"Brackets uses nested context", `person["name"]`, `[{"person":{"name":"Bob"}}]`, "Bob"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := happy.ParserFrom()("KeyName", tt.key)
			checkError(t, err, "could not parse %s into key", tt.key)
			contextJson, err := sample.ParseJson(tt.context)
			checkError(t, err, "could not parse %s into json", tt.context)

			context := happy.ContextOf(contextJson.([]any)...)
			actual := happy.ResolveName(key.(happy.Key), context)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("actual was %s, expected %s", actual, tt.expected)
			}
		})
	}
}

func TestKeyWithInterestingContext(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	context := []any{
		map[string]any{
			"people": []Person{{"Bob", 123}, {"Jim", 456}},
			"type":   func(data any) string { return reflect.TypeOf(data).Elem().Name() },
			"cities": map[string]string{"Europe": "Antwerp", "NorthAmerica": "Chicago"},
		},
		&Person{"Jim", 456},
	}
	tests := []struct {
		name         string
		key          string
		expected     any
		ok           bool
		expectedName string
	}{
		{"At key returns index", `@`, 1, true, "1"},
		{"String key doesn't split on dots", `"people.0"`, "people.0", true, "people.0"},
		{"Function resolves up", `type[.]`, reflect.TypeFor[*Person]().Elem().Name(), true, "Person"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := happy.ParserFrom()("KeyName", tt.key)
			checkError(t, err, "could not parse %s into key", tt.key)

			context := happy.ContextOf(context...)
			actual, ok := key.(happy.Key).Resolve(context)
			actualName := happy.ResolveName(key.(happy.Key), context)

			if ok != tt.ok {
				t.Errorf("ok was %t, expected %t", ok, tt.ok)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("actual was %s, expected %s", actual, tt.expected)
			}
			if !reflect.DeepEqual(actualName, tt.expectedName) {
				t.Errorf("actualName was %s, expected %s", actualName, tt.expectedName)
			}
		})
	}
}
