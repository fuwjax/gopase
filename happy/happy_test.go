package happy_test

import (
	"testing"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/sample"
	"github.com/fuwjax/gopase/when"
)

func TestHappyInterpolation(t *testing.T) {
	render := func(template string, data any) when.WhenOpErr[string] {
		return func() (string, error) {
			return happy.Render(template, data, nil)
		}
	}
	when.YouDoErr("Plaintext", render("This is plain text!", nil)).Expect(t, "This is plain text!")
	when.YouDoErr("Kitty", render("This is (^.^) text!", "kitty")).Expect(t, "This is kitty text!")
	when.YouDoErr("Interpolate", render("This is (^name^) text!", map[string]any{"name": "Bob"})).
		Expect(t, "This is Bob text!")
	when.YouDoErr("Strip leading space", render(`This is 
		
		( ^name^) text!`, map[string]any{"name": "sticky"})).Expect(t, "This issticky text!")
	when.YouDoErr("Strip trailing space", render(`This is    (^name^ ) text!`, map[string]any{"name": "plain"})).
		Expect(t, "This is    plaintext!")
	when.YouDoErr("Handle inner whitespace", render(`This is (^
		name
		^) text!`, map[string]any{"name": "valid"})).Expect(t, "This is valid text!")
}

func TestHappyJsonInterpolation(t *testing.T) {
	render := func(template string, data string) when.WhenOpErr[string] {
		return func() (string, error) {
			json := when.YouErr(sample.ParseJson(data)).ExpectSuccess(t)
			return happy.Render(template, json, nil)
		}
	}
	when.YouDoErr("Simple", render(`
		(^name^)
		(^age^)
		`, `{"name":"Bob","age":123}`)).Expect(t, `
		Bob
		123
		`)
	when.YouDoErr("Dotted", render(`
		(^person.name^)
		(^person.age^)
		`, `{"person":{"name":"Bob Hope","age":55}}`)).Expect(t, `
		Bob Hope
		55
		`)
	when.YouDoErr("Bracketted", render(`
		( ^ person[name] ^ )
		( ^ person[age] ^ )
		`, `{"person":{"name":"Bob Ross","age":35}}`)).Expect(t, `Bob Ross35`)
}
