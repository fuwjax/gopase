package parser

import (
	"fmt"
	"testing"

	"github.com/fuwjax/gopase/funki/testi"
)

func TestBootstrapGramamr(t *testing.T) {
	t.Run("Bootstrap String", func(t *testing.T) {
		fmt.Println(PegGrammar)
	})
}

func TestBootstrapEof(t *testing.T) {
	t.Run("Bootstrap EOF", func(t *testing.T) {
		result, err := Parse("EOF", PegGrammar, PegHandler, "")
		testi.AssertEqual(t, result, "")
		testi.AssertNil(t, err)

		result, err = Parse("EOF", PegGrammar, PegHandler, "a")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at 'a' 1:1 (1) expected not something\nwhile in EOF")
	})
}

func TestBootstrapEol(t *testing.T) {
	t.Run("Bootstrap EOL", func(t *testing.T) {
		result, err := Parse("EOL", PegGrammar, PegHandler, "\n")
		testi.AssertEqual(t, result, "\n")
		testi.AssertNil(t, err)

		result, err = Parse("EOL", PegGrammar, PegHandler, "a")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at 'a' 1:1 (1) expected [\\n\\r]\nwhile in EOL")
	})
}

func TestBootstrapWs(t *testing.T) {
	t.Run("Bootstrap WS", func(t *testing.T) {
		result, err := Parse("WS", PegGrammar, PegHandler, "   ")
		testi.AssertEqual(t, result, "   ")
		testi.AssertNil(t, err)

		result, err = Parse("WS", PegGrammar, PegHandler, " a ")
		testi.AssertEqual(t, result, " ")
		testi.AssertNil(t, err)

		result, err = Parse("WS", PegGrammar, PegHandler, "a ")
		testi.AssertEqual(t, result, "")
		testi.AssertNil(t, err)
	})
}

func TestBootstrapName(t *testing.T) {
	t.Run("Bootstrap Name", func(t *testing.T) {
		result, err := Parse("Name", PegGrammar, PegHandler, "bob")
		testi.AssertEqual(t, result, "bob")
		testi.AssertNil(t, err)

		result, err = Parse("Name", PegGrammar, PegHandler, "B0B ross")
		testi.AssertEqual(t, result, "B0B")
		testi.AssertNil(t, err)

		result, err = Parse("Name", PegGrammar, PegHandler, " bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at ' ' 1:1 (1) expected [_a-zA-Z]\nwhile in Name")

		result, err = Parse("Name", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected [_a-zA-Z]\nwhile in Name")
	})
}

func TestBootstrapPattern(t *testing.T) {
	t.Run("Bootstrap Pattern", func(t *testing.T) {
		result, err := Parse("Pattern", PegGrammar, PegHandler, "[a-z]")
		testi.AssertEqual(t, result, "[a-z]")
		testi.AssertNil(t, err)

		result, err = Parse("Pattern", PegGrammar, PegHandler, "[\\\"]")
		testi.AssertEqual(t, result, "[\\\"]")
		testi.AssertNil(t, err)

		result, err = Parse("Pattern", PegGrammar, PegHandler, "[\"]")
		testi.AssertEqual(t, result, "[\"]")
		testi.AssertNil(t, err)

		result, err = Parse("Pattern", PegGrammar, PegHandler, "[^\\n]")
		testi.AssertEqual(t, result, "[^\\n]")
		testi.AssertNil(t, err)

		result, err = Parse("Pattern", PegGrammar, PegHandler, "[\"]")
		testi.AssertEqual(t, result, "[\"]")
		testi.AssertNil(t, err)

		result, err = Parse("Pattern", PegGrammar, PegHandler, "[]]")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at ']' 1:2 (2) expected \\\nat ']' 1:2 (2) expected [^\\]]\nwhile in Pattern")

		result, err = Parse("Pattern", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected [\nwhile in Pattern")
	})
}

func TestBootstrapComment(t *testing.T) {
	t.Run("Bootstrap Comment", func(t *testing.T) {
		result, err := Parse("Comment", PegGrammar, PegHandler, "#comment")
		testi.AssertEqual(t, result, "#comment")
		testi.AssertNil(t, err)

		result, err = Parse("Comment", PegGrammar, PegHandler, "#")
		testi.AssertEqual(t, result, "#")
		testi.AssertNil(t, err)

		result, err = Parse("Comment", PegGrammar, PegHandler, "not a comment")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at 'n' 1:1 (1) expected #\nwhile in Comment")
	})
}

func TestBootstrapRef(t *testing.T) {
	t.Run("Bootstrap Ref", func(t *testing.T) {
		result, err := Parse("Ref", PegGrammar, PegHandler, "bob")
		testi.AssertEqual(t, result, Ref("bob"))
		testi.AssertNil(t, err)

		result, err = Parse("Ref", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected [_a-zA-Z]\nwhile in Name\nwhile in Ref")
	})
}

func TestBootstrapCharClass(t *testing.T) {
	t.Run("Bootstrap CharClass", func(t *testing.T) {
		result, err := Parse("CharClass", PegGrammar, PegHandler, "[a-z]")
		testi.AssertEqual(t, result, Cls("[a-z]"))
		testi.AssertNil(t, err)

		result, err = Parse("CharClass", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected [\nwhile in Pattern\nwhile in CharClass")
	})
}

func TestBootstrapLiteral(t *testing.T) {
	t.Run("Bootstrap Literal", func(t *testing.T) {
		result, err := Parse("Literal", PegGrammar, PegHandler, "\"hello, world\"")
		testi.AssertEqual(t, result, Lit("hello, world"))
		testi.AssertNil(t, err)

		result, err = Parse("Literal", PegGrammar, PegHandler, "'hello, world'")
		testi.AssertEqual(t, result, Lit("hello, world"))
		testi.AssertNil(t, err)

		result, err = Parse("Literal", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected '\nwhile in SingleLit\nat '1' 1:1 (1) expected \"\nwhile in DoubleLit\nwhile in Literal")
	})
	t.Run("Bootstrap Backslash Primary", func(t *testing.T) {
		result, err := Parse("Literal", PegGrammar, PegHandler, `'\\'`)
		testi.AssertEqual(t, result, Lit(`\`))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapDot(t *testing.T) {
	t.Run("Bootstrap Dot", func(t *testing.T) {
		result, err := Parse("Dot", PegGrammar, PegHandler, ".")
		testi.AssertEqual(t, result, Dot())
		testi.AssertNil(t, err)

		result, err = Parse("Dot", PegGrammar, PegHandler, "1234")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at '1' 1:1 (1) expected .\nwhile in Dot")
	})
}

func TestBootstrapPrimary(t *testing.T) {
	t.Run("Bootstrap Primary", func(t *testing.T) {
		result, err := Parse("Primary", PegGrammar, PegHandler, ".")
		testi.AssertEqual(t, result, Dot())
		testi.AssertNil(t, err)

		result, err = Parse("Primary", PegGrammar, PegHandler, "\"double\"")
		testi.AssertEqual(t, result, Lit("double"))
		testi.AssertNil(t, err)

		result, err = Parse("Primary", PegGrammar, PegHandler, "'single'")
		testi.AssertEqual(t, result, Lit("single"))
		testi.AssertNil(t, err)

		result, err = Parse("Primary", PegGrammar, PegHandler, "[^\"]")
		testi.AssertEqual(t, result, Cls("[^\"]"))
		testi.AssertNil(t, err)

		result, err = Parse("Primary", PegGrammar, PegHandler, "RefName")
		testi.AssertEqual(t, result, Ref("RefName"))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapReqExpr(t *testing.T) {
	t.Run("Bootstrap ReqExpr", func(t *testing.T) {
		result, err := Parse("ReqExpr", PegGrammar, PegHandler, "[0-9]+")
		testi.AssertEqual(t, result, Req(Cls("[0-9]")))
		testi.AssertNil(t, err)

		result, err = Parse("ReqExpr", PegGrammar, PegHandler, "'hi'  +")
		testi.AssertEqual(t, result, Req(Lit("hi")))
		testi.AssertNil(t, err)

		result, err = Parse("ReqExpr", PegGrammar, PegHandler, "Bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at EOF 1:4 (4) expected +\nwhile in ReqExpr")
	})
}

func TestBootstrapRepExpr(t *testing.T) {
	t.Run("Bootstrap RepExpr", func(t *testing.T) {
		result, err := Parse("RepExpr", PegGrammar, PegHandler, "\"yup\"*")
		testi.AssertEqual(t, result, Rep(Lit("yup")))
		testi.AssertNil(t, err)

		result, err = Parse("RepExpr", PegGrammar, PegHandler, ".  *")
		testi.AssertEqual(t, result, Rep(Dot()))
		testi.AssertNil(t, err)

		result, err = Parse("RepExpr", PegGrammar, PegHandler, "Bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at EOF 1:4 (4) expected *\nwhile in RepExpr")
	})
}

func TestBootstrapOptExpr(t *testing.T) {
	t.Run("Bootstrap OptExpr", func(t *testing.T) {
		result, err := Parse("OptExpr", PegGrammar, PegHandler, "RefName?")
		testi.AssertEqual(t, result, Opt(Ref("RefName")))
		testi.AssertNil(t, err)

		result, err = Parse("OptExpr", PegGrammar, PegHandler, ".  ?")
		testi.AssertEqual(t, result, Opt(Dot()))
		testi.AssertNil(t, err)

		result, err = Parse("OptExpr", PegGrammar, PegHandler, "Bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at EOF 1:4 (4) expected ?\nwhile in OptExpr")
	})
}

func TestBootstrapSuffix(t *testing.T) {
	t.Run("Bootstrap Suffix", func(t *testing.T) {
		result, err := Parse("Suffix", PegGrammar, PegHandler, ".+")
		testi.AssertEqual(t, result, Req(Dot()))
		testi.AssertNil(t, err)

		result, err = Parse("Suffix", PegGrammar, PegHandler, "\"double\" *")
		testi.AssertEqual(t, result, Rep(Lit("double")))
		testi.AssertNil(t, err)

		result, err = Parse("Suffix", PegGrammar, PegHandler, "[^\"]?")
		testi.AssertEqual(t, result, Opt(Cls("[^\"]")))
		testi.AssertNil(t, err)

		result, err = Parse("Suffix", PegGrammar, PegHandler, "Bob")
		testi.AssertEqual(t, result, Ref("Bob"))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapNotExpr(t *testing.T) {
	t.Run("Bootstrap NotExpr", func(t *testing.T) {
		result, err := Parse("NotExpr", PegGrammar, PegHandler, "!RefName")
		testi.AssertEqual(t, result, Not(Ref("RefName")))
		testi.AssertNil(t, err)

		result, err = Parse("NotExpr", PegGrammar, PegHandler, "!  .  ?")
		testi.AssertEqual(t, result, Not(Opt(Dot())))
		testi.AssertNil(t, err)

		result, err = Parse("NotExpr", PegGrammar, PegHandler, "Bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at 'B' 1:1 (1) expected !\nwhile in NotExpr")
	})
}

func TestBootstrapAndExpr(t *testing.T) {
	t.Run("Bootstrap AndExpr", func(t *testing.T) {
		result, err := Parse("AndExpr", PegGrammar, PegHandler, "&RefName")
		testi.AssertEqual(t, result, See(Ref("RefName")))
		testi.AssertNil(t, err)

		result, err = Parse("AndExpr", PegGrammar, PegHandler, "&  .  ?")
		testi.AssertEqual(t, result, See(Opt(Dot())))
		testi.AssertNil(t, err)

		result, err = Parse("AndExpr", PegGrammar, PegHandler, "Bob")
		testi.AssertNil(t, result)
		testi.AssertError(t, err, "at 'B' 1:1 (1) expected &\nwhile in AndExpr")
	})
}

func TestBootstrapPrefix(t *testing.T) {
	t.Run("Bootstrap Prefix", func(t *testing.T) {
		result, err := Parse("Prefix", PegGrammar, PegHandler, "!.")
		testi.AssertEqual(t, result, Not(Dot()))
		testi.AssertNil(t, err)

		result, err = Parse("Prefix", PegGrammar, PegHandler, "& \"double\" *")
		testi.AssertEqual(t, result, See(Rep(Lit("double"))))
		testi.AssertNil(t, err)

		result, err = Parse("Prefix", PegGrammar, PegHandler, "[^\"]")
		testi.AssertEqual(t, result, Cls("[^\"]"))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapSeq(t *testing.T) {
	t.Run("Bootstrap Seq", func(t *testing.T) {
		result, err := Parse("Seq", PegGrammar, PegHandler, "A B C")
		testi.AssertEqual(t, result, Seq(Ref("A"), Ref("B"), Ref("C")))
		testi.AssertNil(t, err)

		result, err = Parse("Seq", PegGrammar, PegHandler, ". 'hi' [a-z]")
		testi.AssertEqual(t, result, Seq(Dot(), Lit("hi"), Cls("[a-z]")))
		testi.AssertNil(t, err)

		result, err = Parse("Seq", PegGrammar, PegHandler, "Jim")
		testi.AssertEqual(t, result, Ref("Jim"))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapExpr(t *testing.T) {
	t.Run("Bootstrap Expr", func(t *testing.T) {
		result, err := Parse("Expr", PegGrammar, PegHandler, "A / B / C")
		testi.AssertEqual(t, result, Alt(Ref("A"), Ref("B"), Ref("C")))
		testi.AssertNil(t, err)

		result, err = Parse("Expr", PegGrammar, PegHandler, ". 'hi' / [a-z]")
		testi.AssertEqual(t, result, Alt(Seq(Dot(), Lit("hi")), Cls("[a-z]")))
		testi.AssertNil(t, err)

		result, err = Parse("Expr", PegGrammar, PegHandler, "Jim")
		testi.AssertEqual(t, result, Ref("Jim"))
		testi.AssertNil(t, err)
	})
}

func TestBootstrapParExpr(t *testing.T) {
	t.Run("Bootstrap ParExpr", func(t *testing.T) {
		result, err := Parse("ParExpr", PegGrammar, PegHandler, "(A B C)")
		testi.AssertEqual(t, result, Seq(Ref("A"), Ref("B"), Ref("C")))
		testi.AssertNil(t, err)

		result, err = Parse("ParExpr", PegGrammar, PegHandler, "('hi' / [a-z])")
		testi.AssertEqual(t, result, Alt(Lit("hi"), Cls("[a-z]")))
		testi.AssertNil(t, err)

		result, err = Parse("ParExpr", PegGrammar, PegHandler, "(Jim)")
		testi.AssertEqual(t, result, Ref("Jim"))
		testi.AssertNil(t, err)
	})
}

func TestComplicatedExpr(t *testing.T) {
	t.Run("Bootstrap Complicated", func(t *testing.T) {
		result, err := Parse("Expr", PegGrammar, PegHandler, `'"' (Plain / "\\u" Hex / "\\" Escape)* '"'`)
		testi.AssertEqual(t, result, Seq(Lit(`"`), Rep(Alt(Ref("Plain"), Seq(Lit(`\u`), Ref("Hex")), Seq(Lit(`\`), Ref("Escape")))), Lit(`"`)))
		testi.AssertNil(t, err)
	})
}
