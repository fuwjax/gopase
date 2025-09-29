/*
Package parser handles the general parsing abstractions for gopase.
*/
package parser

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/fuwjax/gopase/funki"
)

type ruleStack struct {
	name string
	next *ruleStack
}

/*
parseCache is a simple named tuple for partial packrat functionality.
*/
type parseCache struct {
	value      any
	err        error
	end        *ParsePosition
	pending    bool
	lrDetected bool
	paths      []string
}

func newCache() *parseCache {
	return &parseCache{nil, errors.New("left recursion detected"), nil, true, false, ([]string)(nil)}
}

/*
ParsePosition is the mark object. There are no public fields or methods on
this type.
*/
type ParsePosition struct {
	grapheme *Grapheme
	cache    map[string]*parseCache
	stack    *ruleStack
	next     *ParsePosition
}

// currently implemented as a linked list to track the current grapheme and
// associated cached Rule results for this position

/*
Creates the initial ParsePostion. Further Positions should be created from
advance().
*/
func newParsePosition(input string) *ParsePosition {
	return &ParsePosition{NewGrapheme(input), make(map[string]*parseCache), nil, nil}
}

/*
Gets a cached result & end mark for a given ref name, if one exists. Return indicates a cache hit.
*/
func (p *ParsePosition) get(name string) (result any, err error, end *ParsePosition, exists bool) {
	cached, exists := p.cache[name]
	if !exists {
		p.stack = &ruleStack{name, p.stack}
		cached = newCache()
		p.cache[name] = cached
	} else if cached.pending {
		cached.lrDetected = true
		for c := p.stack; c.name != name; c = c.next {
			cached.paths = append(cached.paths, c.name)
		}
	}
	return cached.value, cached.err, cached.end, exists
}

/*
Caches a result and end mark for a given ref name. Returns true if ref should recurse.
*/
func (p *ParsePosition) put(name string, result any, err error, end *ParsePosition) (any, error, bool) {
	cached := p.cache[name]
	first := cached.end == nil
	failed := err != nil
	advanced := first || (!failed && cached.end.grapheme.Pos < end.grapheme.Pos)
	detected := advanced && cached.lrDetected
	if advanced {
		cached.value = result
		cached.err = err
		cached.end = end
	}
	if detected {
		for _, n := range cached.paths {
			delete(p.cache, n)
		}
		cached.paths = ([]string)(nil)
		cached.lrDetected = false
	} else {
		p.stack = p.stack.next
		cached.pending = false
	}
	return cached.value, cached.err, cached.pending
}

/*
Advances to next position, creating it if necessary.
*/
func (p *ParsePosition) advance() (*ParsePosition, error) {
	if p.next == nil {
		if p.grapheme.IsEof() {
			return nil, p.grapheme.Error("anything")
		}
		p.next = &ParsePosition{p.grapheme.Next(), make(map[string]*parseCache), nil, nil}
	}
	return p.next, nil
}

/*
Returns an error instance that includes the parse position information.
*/
func (p *ParsePosition) Error(expected string) error {
	return p.grapheme.Error(expected)
}

/*
ParseContext contains the state of a parse.
*/
type ParseContext struct {
	current *ParsePosition
	grammar *Grammar
	handler Handler
}

/*
Create a new ParseContext from the input, rules, and converters.
*/
func newParseContext(input string, grammar *Grammar, handler Handler) *ParseContext {
	return &ParseContext{newParsePosition(input), grammar, handler}
}

/*
Returns a mark that can be passed to Reset() to backtrack
*/
func (c *ParseContext) Mark() *ParsePosition {
	return c.current
}

func (c *ParseContext) At(mark *ParsePosition) bool {
	return c.current == mark
}

/*
Resets to a previous position. Positions must come from Mark() during
this parse operation.
*/
func (c *ParseContext) Reset(mark *ParsePosition) {
	c.current = mark
}

/*
Returns the token at the current parse position.
*/
func (c *ParseContext) Token() string {
	return c.current.grapheme.Token
}

/*
Returns an error instance that includes the parse position information.
*/
func (c *ParseContext) Error(expected string) error {
	return c.current.Error(expected)
}

/*
Returns a substring from the input from the start position to the current position.
*/
func (c *ParseContext) Substring(start *ParsePosition) string {
	var sb strings.Builder
	for p := start; p != c.current; p = p.next {
		sb.WriteString(p.grapheme.Token)
	}
	return sb.String()
}

/*
Accepts the current grapheme and advances to the next one. EOF returns a
non-nil error; other errors may be possible.
*/
func (c *ParseContext) Next() error {
	var err error
	c.current, err = c.current.advance()
	return err
}

/*
Container for reference results while parsing a rule.
*/
type ParseResult struct {
	name  string
	value any
	next  *ParseResult
}

/*
Iterator over name-value pairs.
*/
func (r *ParseResult) Results() iter.Seq2[string, any] {
	return func(yield func(name string, value any) bool) {
		for pr := r; pr != nil; pr = pr.next {
			if !yield(pr.name, pr.value) {
				return
			}
		}
	}
}

/*
Aggregates parse results.
*/
func (r *ParseResult) Chain(sub *ParseResult) *ParseResult {
	if r == nil {
		return sub
	}
	c := r
	for ; c.next != nil; c = c.next {
	}
	c.next = sub
	return r
}

/*
Identifies a right-hand side expression for rule specifications.
*/
type Expr interface {
	fmt.Stringer
	Parse(*ParseContext) (*ParseResult, error)
}

//The end goal is to use a function pointer instead of an interface, right?
//type Expression func(*ParseContext) (*ParseResult, error)

/*
Converts a parse result to an output object.
*/
type Converter func(iter.Seq2[string, any]) (any, error)

/*
Returns a converter for a given rule name.
*/
type Handler func(string) Converter

/*
Uses methods on a type to generate a handler.
*/
func ReflectHandler(handler any) Handler {
	value := reflect.ValueOf(handler)
	return func(name string) Converter {
		method := value.MethodByName(name)
		if !method.IsValid() {
			return nil
		}
		return func(result iter.Seq2[string, any]) (any, error) {
			methodArgs := make([]reflect.Value, 1)
			methodArgs[0] = reflect.ValueOf(result)
			returns := method.Call(methodArgs)
			err := returns[1].Interface()
			if err == nil {
				return returns[0].Interface(), nil
			}
			return returns[0].Interface(), err.(error)
		}
	}
}

/*
Wraps a map, struct, or nil into a valid handler.
*/
func WrapHandler(handler any) Handler {
	if handler == nil {
		return func(string) Converter {
			return nil
		}
	}
	h, ok := handler.(Handler)
	if ok {
		return h
	}
	m, ok := handler.(map[string]Converter)
	if ok {
		return func(key string) Converter {
			return m[key]
		}
	}
	return ReflectHandler(handler)
}

/*
Defines a grammar Rule.
*/
type Rule struct {
	name string
	expr Expr
}

/*
Creates a rule.
*/
func NewRule(name string, expr Expr) *Rule {
	return &Rule{name, expr}
}

/*
Parses the input and returns a converted output object.
*/
func (r *Rule) Parse(context *ParseContext) (any, error) {
	var mark *ParsePosition
	if context == nil || r == nil {
		return nil, fmt.Errorf("WTF")
	}
	converter := context.handler(r.name)
	if converter == nil { // just to avoid the reference count
		mark = context.Mark()
	}
	result, err := r.expr.Parse(context)
	if err != nil {
		return nil, fmt.Errorf("%s\nwhile in %s", err, r.name)
	}
	if converter != nil {
		return converter(result.Results())
	}
	return context.Substring(mark), nil
}

func (r *Rule) String() string {
	return fmt.Sprintf("Rule(\"%s\", %s)", r.name, r.expr)
}

/*
Represents the collection of rules that specifies a grammar.
*/
type Grammar struct {
	rules map[string]*Rule
	order []string
}

/*
Creates an empty grammar.
*/
func NewGrammar() *Grammar {
	return &Grammar{make(map[string]*Rule), make([]string, 0)}
}

/*
Creates and adds a rule to the grammar.
*/
func (g *Grammar) AddRule(name string, expr Expr) *Grammar {
	return g.Add(NewRule(name, expr))
}

/*
Adds a rule to the grammar.
*/
func (g *Grammar) Add(rule *Rule) *Grammar {
	g.rules[rule.name] = rule
	g.order = append(g.order, rule.name)
	return g
}

/*
Returns a rule by name.
*/
func (g *Grammar) Rule(name string) *Rule {
	return g.rules[name]
}

/*
Returns the rules specifying this grammar.
*/
func (g *Grammar) Rules() iter.Seq2[string, *Rule] {
	return func(yield func(string, *Rule) bool) {
		for _, name := range g.order {
			if !yield(name, g.rules[name]) {
				return
			}
		}
	}
}

func (g *Grammar) String() string {
	return strings.Join(funki.Apply(slices.Collect(funki.Values(g.Rules())), (*Rule).String), "\n")
}

type Parser[T any] func(input string) (T, error)

type ParserFrom func(root, input string) (any, error)

/*
Creates a new parser.
*/
func NewParser[T any](root string, grammar string, handler any) Parser[T] {
	parser := NewParserFrom(grammar, handler)
	return func(input string) (T, error) {
		result, err := parser(root, input)
		if err != nil {
			var t T
			return t, err
		}
		if result == nil {
			var zero T
			return zero, nil
		}
		return result.(T), nil
	}
}

func NewParserFrom(grammar string, handler any) ParserFrom {
	rules, err := Bootstrap(grammar)
	realHandler := WrapHandler(handler)
	return func(root, input string) (any, error) {
		if err != nil {
			return nil, err
		}
		return Parse(root, rules, realHandler, input)
	}
}

func BootstrapParser[T any](root string, grammar *Grammar, handler Handler) Parser[T] {
	return func(input string) (T, error) {
		result, err := Parse(root, grammar, handler, input)
		if err != nil {
			var t T
			return t, err
		}
		return result.(T), nil
	}
}

func BootstrapParserFrom(grammar *Grammar, handler Handler) ParserFrom {
	return func(root, input string) (any, error) {
		return Parse(root, grammar, handler, input)
	}
}

/*
Parses the input according to the root, grammar, and handler.
*/
func Parse(root string, grammar *Grammar, handler Handler, input string) (any, error) {
	ref := Ref(root)
	result, err := ref.Parse(newParseContext(input, grammar, handler))
	if err != nil {
		return nil, err
	}
	return result.value, nil
}
