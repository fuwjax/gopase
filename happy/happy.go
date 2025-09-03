package happy

import (
	"fmt"
	"strings"
)

func Render(template string, data any, partials map[string]*Template) (string, error) {
	comp, err := Compile(template)
	if err != nil {
		return "", err
	}
	return comp.Render(ContextOf(data), partials)
}

type Context struct {
	Index any
	Data  any
	Next  *Context
}

func ContextOf(data ...any) *Context {
	var c *Context
	for i, d := range data {
		c = c.With(i, d)
	}
	return c
}

func (c *Context) With(index, data any) *Context {
	return &Context{index, data, c}
}

type Renderer interface {
	Render(context *Context, partials map[string]*Template) (string, error)
}

type Template struct {
	Content []Renderer
}

func (t *Template) Render(context *Context, partials map[string]*Template) (string, error) {
	if partials == nil {
		partials = make(map[string]*Template)
	}
	var sb strings.Builder
	for _, snippet := range t.Content {
		text, err := snippet.Render(context, partials)
		if err != nil {
			return "", err
		}
		sb.WriteString(text)
	}
	return sb.String(), nil
}

type Section struct {
	Name    Key
	Content *Template
}

func (s *Section) Render(context *Context, partials map[string]*Template) (string, error) {
	data, ok := s.Name.Resolve(context)
	if data == nil || !ok {
		return "", nil
	}
	slice, ok := Iter(data)
	if !ok {
		return s.Content.Render(context.With(nil, data), partials)
	}
	var sb strings.Builder
	for index, data := range slice {
		result, err := s.Content.Render(context.With(index, data), partials)
		if err != nil {
			return "", err
		}
		sb.WriteString(result)
	}
	return sb.String(), nil
}

type Reference struct {
	Name Key
}

func (r *Reference) Render(context *Context, partials map[string]*Template) (string, error) {
	data, ok := r.Name.Resolve(context)
	if !ok && data == nil {
		return "", nil
	}
	str, ok := data.(string)
	if ok {
		return str, nil
	}
	return fmt.Sprint(data), nil
}

type Plaintext struct {
	Text string
}

func (p *Plaintext) Render(context *Context, partials map[string]*Template) (string, error) {
	return p.Text, nil
}

type Include struct {
	Name Key
}

func (i *Include) Render(context *Context, partials map[string]*Template) (string, error) {
	name := ResolveName(i.Name, context)
	partial, ok := partials[name]
	if !ok {
		return "", fmt.Errorf("no partial named %s", name)
	}
	return partial.Render(context, partials)
}

type Partial struct {
	Name    Key
	Content *Template
}

func (p *Partial) Render(context *Context, partials map[string]*Template) (string, error) {
	name := ResolveName(p.Name, context)
	if name == "" {
		return "", fmt.Errorf("partial cannot be given empty name")
	}
	partials[name] = p.Content
	return "", nil
}
