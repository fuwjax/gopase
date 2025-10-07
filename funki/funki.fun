Library r -> $NewLibrary(r.[]Declaration)
Declaration r -> $NewFunction(r.Name, r.Relation)
Relation r -> test r
    {Production: production} -> $NewRelation(r.Pattern, production)
    {Relation: relation} -> $NewRelation(r.Pattern, relation)
Pattern r -> r._
ListPattern r -> $NewListPattern(r.[]Pattern, r.Name)
TuplePattern r -> $NewTuplePattern(r.[]Pattern)
ObjectPattern r -> $NewObjectPattern(r.[]EntryPattern)
EntryPattern r -> (r.Name, r.Pattern)
RelationPattern r -> $NewRelationPattern(r.Name, r.Pattern)

Production r -> r._
TestProduction r -> $NewTest(r.Production, r.[]Relation)
Invocation r -> $NewInvoke(r.Name, r.Production)
ObjectProduction r -> $NewList(r.[]EntryProduction)
EntryProduction r -> (r.Name, r.Production)
TupleProduction r -> $NewTuple(r.[]Production)
ListProduction r -> $NewList(r.[]Production)

MayResolveToObject r -> r._
MayResolveToList r -> r._
MayResolveToNumber r -> r._

Operation r -> r._
ObjectOperation r -> r._
Dot r -> $Get(r.MayResolveToObject, r.Name)
DotUnder r -> $GetFirst(r.MayResolveToObject)
DotList r -> $GetList(r.MayResolveToObject, r.Name)
Concatenate r-> $Concatenate(r.[]MayResolveToList)
Contains r -> $Contains(r.MayResolveToList, r.Production)
Equality [(_, x):(_, op):(_, y)] -> $Infix(op, x, y)
Comparison [(_, x):(_, op):(_, y)] -> $Infix(op, x, y)
Calculation [(_, x):(_, op):(_, y)] -> $Infix(op, x, y)
PrefixOperation [(_, op):(_, x)] -> $Prefix(op, x)

    