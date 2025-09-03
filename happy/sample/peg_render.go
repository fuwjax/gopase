package sample

import (
	"reflect"
	"sync"

	"github.com/fuwjax/gopase/happy"
	"github.com/fuwjax/gopase/parser"
)

const template = `
( ^=Sequence^)parser.Seq((^*exprs^)(^*@^), (^/^)(^>type[.]^)(^/^))(^/^)
( ^=Options^)parser.Alt((^*exprs^)(^*@^), (^/^)(^>type[.]^)(^/^))(^/^)
( ^=Optional^)parser.Opt((^*expr^)(^>type[.]^)(^/^))(^/^)
( ^=Repeated^)parser.Rep((^*expr^)(^>type[.]^)(^/^))(^/^)
( ^=Required^)parser.Req((^*expr^)(^>type[.]^)(^/^))(^/^)
( ^=CharClass^)parser.Cls("(^regex^)")(^/^)
( ^=Literal^)parser.Lit("(^literal^)")(^/^)
( ^=Any^)parser.Dot()(^/^)
( ^=Reference^)parser.Ref("(^name^)")(^/^)
( ^=PositiveLookahead^)parser.See((^*expr^)(^>type[.]^)(^/^))(^/^)
( ^=NegativeLookahead^)parser.Not((^*expr^)(^>type[.]^)(^/^))(^/^ )

package (^package^)

import "github.com/fuwjax/gopase/parser"

func (^name^)Grammar() *Grammar {
	grammar := parser.NewGrammar()
	(^*grammar.Rules^ )
	grammar.AddRule("(^@^)", (^*expr^)(^>type[.]^)(^/^))
	(^/^ )
	return grammar
}
`

var PegTemplate = sync.OnceValues(func() (*happy.Template, error) {
	return happy.Compile(template)
})

func typeOf(data any) string {
	return reflect.TypeOf(data).Elem().Name()
}

func RenderPeg(pack, name string, grammar *parser.Grammar) (string, error) {
	template, err := PegTemplate()
	if err != nil {
		return "", err
	}
	context := happy.ContextOf(map[string]any{"package": pack, "name": name, "grammar": grammar, "type": typeOf})
	return template.Render(context, nil)
}
