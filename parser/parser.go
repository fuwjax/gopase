/*
Package parser handles the general parsing abstractions for gopase.
*/
package parser

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"
)

/*
parseCache is a simple named tuple for partial packrat functionality.
*/
type parseCache struct {
	value any
	end   *ParsePosition
}

/*
ParsePosition is the mark object. There are no public fields or methods on
this type.
*/
type ParsePosition struct {
	grapheme *Grapheme
	cache    map[string]*parseCache
	next     *ParsePosition
}

// currently implemented as a linked list to track the current grapheme and
// associated cached Rule results for this position

/*
Creates the initial ParsePostion. Further Positions should be created from
advance().
*/
func newParsePosition(input string) *ParsePosition {
	return &ParsePosition{NewGrapheme(input), nil, nil}
}

/*
Gets a cached result & end mark for a given ref name, if one exists.
*/
func (p *ParsePosition) get(name string) (result any, end *ParsePosition, exists bool) {
	if p.cache == nil {
		return nil, nil, false
	}
	cached, exists := p.cache[name]
	if !exists {
		return nil, nil, false
	}
	return cached.value, cached.end, exists
}

/*
Caches a result and end mark for a given ref name.
*/
func (p *ParsePosition) put(name string, result any, end *ParsePosition) {
	if p.cache == nil {
		p.cache = make(map[string]*parseCache)
	}
	p.cache[name] = &parseCache{result, end}
}

/*
Advances to next position, creating it if necessary.
*/
func (p *ParsePosition) advance() (*ParsePosition, error) {
	if p.next == nil {
		if p.grapheme.IsEof() {
			return nil, p.grapheme.Error("anything")
		}
		p.next = &ParsePosition{p.grapheme.Next(), nil, nil}
	}
	return p.next, nil
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
	return c.current.grapheme.Error(expected)
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
type Converter func(*ParseResult) (any, error)

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
		return func(result *ParseResult) (any, error) {
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
		return Map2Func(m)
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
	converter := context.handler(r.name)
	if converter == nil { // just to avoid the reference count
		mark = context.Mark()
	}
	result, err := r.expr.Parse(context)
	if err != nil {
		return nil, err
	}
	if converter != nil {
		return converter(result)
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
func (g *Grammar) Rules() iter.Seq[*Rule] {
	return func(yield func(*Rule) bool) {
		for _, name := range g.order {
			if !yield(g.rules[name]) {
				return
			}
		}
	}
}

/*
Creates a new parser from this grammar with the given root and handler.
*/
func (g *Grammar) Parser(root string, handler any) *Parser {
	return NewParser(root, handler, g)
}

func (g *Grammar) String() string {
	return strings.Join(Apply(slices.Collect(g.Rules()), (*Rule).String), "\n")
}

/*
Encapsulates the shared portion. Has the same sharing semantics as the handler.
*/
type Parser struct {
	grammar *Grammar
	root    string
	handler Handler
}

/*
Creates a new parser.
*/
func NewParser(root string, handler any, grammar *Grammar) *Parser {
	return &Parser{grammar, root, WrapHandler(handler)}
}

/*
Parses the input according to the grammar and handler for this parser.
*/
func (p *Parser) Parse(input string) (any, error) {
	return Parse(p.root, p.grammar, p.handler, input)
}

/*
Parses the input from the new root rule.
*/
func (p *Parser) ParseFrom(root string, input string) (any, error) {
	return Parse(root, p.grammar, p.handler, input)
}

func (p *Parser) String() string {
	return fmt.Sprintf("# %s\n%s", p.root, p.grammar)
}

/*
Parses the input according to the root, grammar, and handler.
*/
func Parse(root string, grammar *Grammar, handler Handler, input string) (any, error) {
	return grammar.Rule(root).Parse(newParseContext(input, grammar, handler))
}
