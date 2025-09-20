/*
Happy is templating for the soul. It's intended to be visually pleasing with tags that
look like little emoji faces, and designed to be as simple as possible while still
retaining the power required for meaningful general purpose templating use cases.

The main way of interacting with happy is via Render

	output, err := happy.Render(aTemplateString, someData, anyPartials)

This will parse the aTemplateString into a happy.Template, convert someData into a happy.Context and pass this context and anyParials
to template.Render. Which makes the Render function equivalent to.

	template, err := happy.Compile(aTemplateString)
	if err == nil {
		context := happy.ContextOf(data)
		output, err := template.Render(context, anyPartials)
	}
*/
package happy

import (
	"fmt"
	"maps"
	"strings"
)

//Partials and data may be nil, but it'll be a pretty boring render.
/*
Renders data against a template with supporting partials.
*/
func Render(template string, data any, partials map[string]Template) (string, error) {
	comp, err := Compile(template)
	if err != nil {
		return "", err
	}
	return comp.Render(ContextOf(data), partials)
}

/*
Represents a pattern for applying the context to produce some output.
*/
type Template interface {
	Render(context Context, partials map[string]Template) (string, error)
}

// The following are constructors made public for testing.

func Content(content []Template) Template {
	return &template{content}
}

func Section(name Key, content Template) Template {
	return &section{name, content}
}

func Invert(name Key, content Template) Template {
	return &invert{name, content}
}

func Reference(name Key) Template {
	return &reference{name}
}

func Plaintext(text string) Template {
	return &plaintext{text}
}

func Include(name Key) Template {
	return &include{name}
}

func Partial(name Key, content Template) Template {
	return &partial{name, content}
}

type template struct {
	content []Template
}

func (t *template) Render(context Context, partials map[string]Template) (string, error) {
	if partials == nil {
		partials = make(map[string]Template)
	} else {
		partials = maps.Clone(partials)
	}
	var sb strings.Builder
	for _, snippet := range t.content {
		text, err := snippet.Render(context, partials)
		if err != nil {
			return "", err
		}
		sb.WriteString(text)
	}
	return sb.String(), nil
}

type section struct {
	name    Key
	content Template
}

func (s *section) Render(context Context, partials map[string]Template) (string, error) {
	data, ok := s.name.Resolve(context)
	if !ok || data == nil {
		return "", nil
	}
	slice, ok := Iter(data)
	if !ok {
		if !Truthy(data) {
			return "", nil
		}
		return s.content.Render(context.With(nil, data), partials)
	}
	var sb strings.Builder
	for index, data := range slice {
		result, err := s.content.Render(context.With(index, data), partials)
		if err != nil {
			return "", err
		}
		sb.WriteString(result)
	}
	return sb.String(), nil
}

type invert struct {
	name    Key
	content Template
}

func (s *invert) Render(context Context, partials map[string]Template) (string, error) {
	data, ok := s.name.Resolve(context)
	if ok && Truthy(data) {
		return "", nil
	}
	return s.content.Render(context, partials)
}

type reference struct {
	name Key
}

func (r *reference) Render(context Context, partials map[string]Template) (string, error) {
	data, ok := r.name.Resolve(context)
	if !ok && !Truthy(data) {
		return "", nil
	}
	str, ok := data.(string)
	if ok {
		return str, nil
	}
	return fmt.Sprint(data), nil
}

type plaintext struct {
	text string
}

func (p *plaintext) Render(context Context, partials map[string]Template) (string, error) {
	return p.text, nil
}

type include struct {
	name Key
}

func (i *include) Render(context Context, partials map[string]Template) (string, error) {
	name := ResolveName(i.name, context)
	partial, ok := partials[name]
	if !ok {
		return "", fmt.Errorf("no partial named %s", name)
	}
	return partial.Render(context, partials)
}

type partial struct {
	name    Key
	content Template
}

func (p *partial) Render(context Context, partials map[string]Template) (string, error) {
	name := ResolveName(p.name, context)
	if name == "" {
		return "", fmt.Errorf("partial cannot be given empty name")
	}
	partials[name] = p.content
	return "", nil
}
