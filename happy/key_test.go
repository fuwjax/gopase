package happy_test

import (
	"reflect"
	"testing"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/sample"
	"github.com/fuwjax/gopase/when"
)

func MatchJson(expected string) when.Matcher[any] {
	return func(t *testing.T, actual any) bool {
		json := when.YouErr(sample.ParseJson(expected)).ExpectSuccess(t)
		return when.AssertEqual(t, actual, json)
	}
}

func TestKeyResolve(t *testing.T) {
	resolve := func(key string, context string) when.WhenOpOk[any] {
		return func() (any, bool) {
			key := when.YouErr(happy.ParserFrom()("KeyName", key)).ExpectSuccess(t)
			contextJson := when.YouErr(sample.ParseJson(context)).ExpectSuccess(t)
			context := happy.ContextOf(contextJson.([]any)...)
			return key.(happy.Key).Resolve(context)
		}
	}
	when.YouDoOk("String key", resolve(`"name"`, `[]`)).ExpectMatch(t, MatchJson(`"name"`))
	when.YouDoOk("String key ignores context", resolve(`"name"`, `[{"name":"Bob"}]`)).ExpectMatch(t, MatchJson(`"name"`))
	when.YouDoOk("Plain key uses context", resolve(`name`, `[{"name":"Bob"}]`)).ExpectMatch(t, MatchJson(`"Bob"`))
	when.YouDoOk("Plain key fails without context", resolve(`name`, `[{"first_name":"Bob"}]`)).ExpectFailure(t)
	when.YouDoOk("Plain key uses latest context", resolve(`name`, `[{"name":"Bob"},{"name":"Jim"}]`)).ExpectMatch(t, MatchJson(`"Jim"`))
	when.YouDoOk("Plain key travels context", resolve(`name`, `[{"name":"Bob"},{"first_name":"Jim"}]`)).ExpectMatch(t, MatchJson(`"Bob"`))
	when.YouDoOk("Dot returns context", resolve(`.`, `[{"name":"Bob"}]`)).ExpectMatch(t, MatchJson(`{"name":"Bob"}`))
	when.YouDoOk("Dotted uses nested context", resolve(`person.name`, `[{"person":{"name":"Bob"}}]`)).ExpectMatch(t, MatchJson(`"Bob"`))
	when.YouDoOk("Brackets uses nested context", resolve(`person["name"]`, `[{"person":{"name":"Bob"}}]`)).ExpectMatch(t, MatchJson(`"Bob"`))
}

func TestKeyResolveName(t *testing.T) {
	resolveName := func(tag string, data string) when.WhenOp[string] {
		return func() string {
			key := when.YouErr(happy.ParserFrom()("KeyName", tag)).ExpectSuccess(t)
			contextJson := when.YouErr(sample.ParseJson(data)).ExpectSuccess(t)
			context := happy.ContextOf(contextJson.([]any)...)
			return happy.ResolveName(key.(happy.Key), context)
		}
	}
	when.YouDo("String key", resolveName(`"name"`, `[]`)).Expect(t, "name")
	when.YouDo("String key ignores context", resolveName(`"name"`, `[{"name":"Bob"}]`)).Expect(t, "name")
	when.YouDo("Plain key uses context", resolveName(`name`, `[{"name":"Bob"}]`)).Expect(t, "Bob")
	when.YouDo("Plain key reverts to literal without context", resolveName(`name`, `[{"first_name":"Bob"}]`)).Expect(t, "name")
	when.YouDo("Plain key uses latest context", resolveName(`name`, `[{"name":"Bob"},{"name":"Jim"}]`)).Expect(t, "Jim")
	when.YouDo("Plain key travels context", resolveName(`name`, `[{"name":"Bob"},{"first_name":"Jim"}]`)).Expect(t, "Bob")
	when.YouDo("Dot returns String() of current context", resolveName(`.`, `[{"name":"Bob"}]`)).Expect(t, "map[name:Bob]")
	when.YouDo("Dotted uses nested context", resolveName(`person.name`, `[{"person":{"name":"Bob"}}]`)).Expect(t, "Bob")
	when.YouDo("Brackets uses nested context", resolveName(`person["name"]`, `[{"person":{"name":"Bob"}}]`)).Expect(t, "Bob")
}

func TestKeyWithInterestingContext(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	context := happy.ContextOf(
		map[string]any{
			"people": []Person{{"Bob", 123}, {"Jim", 456}},
			"type":   func(data any) string { return reflect.TypeOf(data).Elem().Name() },
			"cities": map[string]string{"Europe": "Antwerp", "NorthAmerica": "Chicago"},
		},
		&Person{"Jim", 456})
	resolve := func(key string) when.WhenOpOk[any] {
		return func() (any, bool) {
			key := when.YouErr(happy.ParserFrom()("KeyName", key)).ExpectSuccess(t)
			return key.(happy.Key).Resolve(context)
		}
	}
	resolveName := func(tag string) when.WhenOp[string] {
		return func() string {
			key := when.YouErr(happy.ParserFrom()("KeyName", tag)).ExpectSuccess(t)
			return happy.ResolveName(key.(happy.Key), context)
		}
	}

	when.YouDoOk("At key returns index", resolve(`@`)).Expect(t, 1)
	when.YouDo("At key returns index", resolveName(`@`)).Expect(t, "1")

	when.YouDoOk("String key doesn't split on dots", resolve(`"people.0"`)).Expect(t, "people.0")
	when.YouDo("String key doesn't split on dots", resolveName(`"people.0"`)).Expect(t, "people.0")

	when.YouDoOk("Function resolves up", resolve(`type[.]`)).Expect(t, reflect.TypeFor[*Person]().Elem().Name())
	when.YouDo("Function resolves up", resolveName(`type[.]`)).Expect(t, "Person")
}
