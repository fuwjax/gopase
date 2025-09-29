Grammar r -> {rules: FilterNil r.[]Line}
Line [p:_] -> test p
	"Rule" r -> r
	_ -> nil
Rule r -> {name: r.Name, expr: r.Expr}
Expr r -> {exprs: r.[]Seq}
Seq r -> {exprs: r.[]Prefix}
Prefix r -> r._
AndExpr r -> {expr: r.Suffix}
NotExpr r -> {expr: r.Suffix}
Suffix r -> r._
OptExpr r -> {expr: r.Primary}
RepExpr r -> {expr: r.Primary}
ReqExpr r -> {expr: r.Primary}
Primary r -> r._
Dot _ -> {}
ParExpr r -> r.Expr
Literal r -> r._
CharClass results -> {regex: results.Pattern}
Ref results -> {name: results.Name}
DoubleLit result -> test result
	[] -> ""
	[("DoubleEscape" match):pairs] -> Escape match ++ DoubleLit(pairs)
	[("DoublePlain" match):pairs] -> match ++ DoubleLit(pairs)

Escape ch -> test ch
	"n" -> "\n"
	"r" -> "\r"
	"t" -> "\t"
	x -> x
FilterNil list -> test list
	[] -> []
	[nil:ps] -> FilterNil ps
	[e:ps] -> e : FilterNil ps
