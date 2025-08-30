package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Sequence struct {
	exprs []Expr
}

func Seq(exprs ...Expr) Expr {
	if len(exprs) == 1 {
		return exprs[0]
	}
	return &Sequence{exprs}
}

func (x *Sequence) Parse(context *ParseContext) (*ParseResult, error) {
	var result *ParseResult
	for _, expr := range x.exprs {
		res, err := expr.Parse(context)
		if err != nil {
			return nil, err
		}
		result = result.Chain(res)
	}
	return result, nil
}

func (x *Sequence) String() string {
	return "Seq(" + strings.Join(Apply(x.exprs, Expr.String), ", ") + ")"
}

type Options struct {
	exprs []Expr
}

func Alt(exprs ...Expr) Expr {
	if len(exprs) == 1 {
		return exprs[0]
	}
	return &Options{exprs}
}

func (x *Options) Parse(context *ParseContext) (*ParseResult, error) {
	mark := context.Mark()
	var poly PolyError
	for _, expr := range x.exprs {
		result, err := expr.Parse(context)
		if err == nil {
			return result, nil
		}
		poly.Add(err)
		context.Reset(mark)
	}
	return nil, &poly
}

func (x *Options) String() string {
	return "Alt(" + strings.Join(Apply(x.exprs, Expr.String), ",") + ")"
}

type Optional struct {
	expr Expr
}

func Opt(expr Expr) Expr {
	return &Optional{expr}
}

func (x *Optional) Parse(context *ParseContext) (*ParseResult, error) {
	mark := context.Mark()
	result, err := x.expr.Parse(context)
	if err != nil {
		context.Reset(mark)
		return nil, nil
	}
	return result, nil
}

func (x *Optional) String() string {
	return fmt.Sprintf("Opt(%s)", x.expr)
}

type Repeated struct {
	expr Expr
}

func Rep(expr Expr) Expr {
	return &Repeated{expr}
}

func (x *Repeated) Parse(context *ParseContext) (*ParseResult, error) {
	var agg *ParseResult
	for {
		mark := context.Mark()
		result, err := x.expr.Parse(context)
		if err != nil || context.At(mark) {
			context.Reset(mark)
			break
		}
		agg = agg.Chain(result)
	}
	return agg, nil
}

func (x *Repeated) String() string {
	return fmt.Sprintf("Rep(%s)", x.expr)
}

type Required struct {
	expr Expr
}

func Req(expr Expr) Expr {
	return &Required{expr}
}

func (x *Required) Parse(context *ParseContext) (*ParseResult, error) {
	var agg *ParseResult
	result, err := x.expr.Parse(context)
	if err != nil {
		return nil, err
	}
	agg = agg.Chain(result)
	for {
		mark := context.Mark()
		result, err = x.expr.Parse(context)
		if err != nil || context.At(mark) {
			context.Reset(mark)
			break
		}
		agg = agg.Chain(result)
	}
	return agg, nil
}

func (x *Required) String() string {
	return fmt.Sprintf("Req(%s)", x.expr)
}

type CharClass struct {
	regex *regexp.Regexp
}

func Cls(pattern string) Expr {
	return &CharClass{regexp.MustCompile(pattern)}
}

func (x *CharClass) Parse(context *ParseContext) (*ParseResult, error) {
	if !x.regex.MatchString(context.Token()) {
		return nil, context.Error(x.regex.String())
	}
	err := context.Next()
	return nil, err
}

func (x *CharClass) String() string {
	return fmt.Sprintf("Cls(\"%s\")", x.regex)
}

type Literal struct {
	literal string
}

func Lit(literal string) Expr {
	return &Literal{literal}
}

func (x *Literal) Parse(context *ParseContext) (*ParseResult, error) {
	for ch := range Graphemes(x.literal) {
		if ch.Token != context.Token() {
			return nil, context.Error(ch.Token)
		}
		err := context.Next()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (x *Literal) String() string {
	return fmt.Sprintf("Lit(`%s`)", x.literal)
}

type Any struct{}

func Dot() Expr {
	return &Any{}
}

func (x *Any) Parse(context *ParseContext) (*ParseResult, error) {
	err := context.Next()
	return nil, err
}

func (x *Any) String() string {
	return "Dot()"
}

type Reference struct {
	name string
}

func Ref(name string) Expr {
	return &Reference{name}
}

func (x *Reference) Parse(context *ParseContext) (*ParseResult, error) {
	mark := context.Mark()
	var result any
	var err error
	result, end, ok := mark.get(x.name)
	if ok {
		context.Reset(end)
	} else {
		rule := context.grammar.Rule(x.name)
		if rule == nil {
			return nil, fmt.Errorf("no such rule: %s", x.name)
		}
		result, err = rule.Parse(context)
		if err != nil {
			return nil, err
		}
		mark.put(x.name, result, context.Mark())
	}
	return &ParseResult{x.name, result, nil}, nil
}

func (x *Reference) String() string {
	return fmt.Sprintf("Ref(\"%s\")", x.name)
}

type PositiveLookahead struct {
	expr Expr
}

func See(expr Expr) Expr {
	return &PositiveLookahead{expr}
}

func (x *PositiveLookahead) Parse(context *ParseContext) (*ParseResult, error) {
	mark := context.Mark()
	_, err := x.expr.Parse(context)
	context.Reset(mark) // forces zero length, but only meaningful after a match
	return nil, err
}

func (x *PositiveLookahead) String() string {
	return fmt.Sprintf("See(%s)", x.expr)
}

type NegativeLookahead struct {
	expr Expr
}

func Not(expr Expr) Expr {
	return &NegativeLookahead{expr}
}

func (x *NegativeLookahead) Parse(context *ParseContext) (*ParseResult, error) {
	mark := context.Mark()
	_, err := x.expr.Parse(context)
	context.Reset(mark)
	if err != nil {
		return nil, nil
	}
	return nil, mark.grapheme.Error("not something")
}

func (x *NegativeLookahead) String() string {
	return fmt.Sprintf("Not(%s)", x.expr)
}
