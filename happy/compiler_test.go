package happy_test

import (
	"testing"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/when"
)

func TestCompileKeyName(t *testing.T) {
	parseKeyName := func(input string) when.WhenOpErr[any] {
		return func() (any, error) {
			return happy.ParserFrom()("KeyName", input)
		}
	}
	when.YouDoErr("String key", parseKeyName(`"name"`)).Expect(t, happy.Lit("name"))
	when.YouDoErr("Plain key", parseKeyName(`name`)).Expect(t, happy.Bracket("name"))
	when.YouDoErr("Dotted key", parseKeyName(`person.name`)).
		Expect(t, happy.Dotted(happy.Bracket("person"), happy.Bracket("name")))
	when.YouDoErr("Bracket key", parseKeyName(`person[name]`)).
		Expect(t, happy.Bracket("person", happy.Bracket("name")))
	when.YouDoErr("Crazy key", parseKeyName(`cities["Europe"]["Amsterdam"].people.0.address[division.name].upper[.]`)).
		Expect(t, happy.Dotted(
			happy.Bracket("cities", happy.Lit("Europe"), happy.Lit("Amsterdam")),
			happy.Bracket("people"),
			happy.Bracket("0"),
			happy.Bracket("address", happy.Dotted(happy.Bracket("division"), happy.Bracket("name"))),
			happy.Bracket("upper", happy.Dot())))
}

func TestCompileTags(t *testing.T) {
	parseTag := func(input string) when.WhenOpErr[any] {
		return func() (any, error) {
			return happy.ParserFrom()("Tag", input)
		}
	}
	when.YouDoErr("Section", parseTag(`(^*name^)(^.^)(^/^)`)).
		Expect(t, &happy.Section{happy.Bracket("name"), &happy.Template{[]happy.Renderer{&happy.Reference{happy.Dot()}}}})
	when.YouDoErr("Comment", parseTag(`(^# a comment^)`)).
		Expect(t, nil)
	when.YouDoErr("Value", parseTag(`(^name^)`)).
		Expect(t, &happy.Reference{happy.Bracket("name")})
	when.YouDoErr("Value consumes ws", parseTag(`    ( ^name^ )      `)).
		Expect(t, &happy.Reference{happy.Bracket("name")})
	when.YouDoErr("Else", parseTag(`(^!name^)`)).
		Expect(t, &happy.Reference{happy.Invert(happy.Bracket("name"))})
	when.YouDoErr("Include", parseTag(`(^>name^)`)).
		Expect(t, &happy.Include{happy.Bracket("name")})
}
