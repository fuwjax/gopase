package funki

type Library struct {
	functions []Function
}

func NewLibrary(functions []Function) *Library {
	return &Library{functions}
}

type Function struct {
	name     string
	relation Relation
}

func NewFunction(name string, relation Relation) *Function {
	return &Function{name, relation}
}

type Relation struct {
	pattern    Pattern
	production Production
}

func NewRelation(pattern Pattern, production Production) *Relation {
	return &Relation{pattern, production}
}

type Pattern interface {
}

type Production interface {
}

type Underbar struct{}

type Reference struct {
	name string
}

func NewReference(name string) *Reference {
	return &Reference{name}
}

type Exact struct {
	value any
}

func NewExact(value any) *Exact {
	return &Exact{value}
}
