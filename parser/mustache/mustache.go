package mustache

import (
	"fmt"
	"strings"
)

func Render(template string, data any) (string, error) {
	comp, err := Compile(template)
	if err != nil {
		return "", err
	}
	return comp.Render(data)
}

type Renderer interface {
	Render(data any) (string, error)
}

type Template struct {
	snippets []Renderer
}

func (t *Template) Render(data any) (string, error) {
	var sb strings.Builder
	for _, snippet := range t.snippets {
		text, err := snippet.Render(data)
		if err != nil {
			return "", err
		}
		sb.WriteString(text)
	}
	return sb.String(), nil
}

type Section struct {
	Name  string
	Inner Renderer
}

func (s *Section) Render(data any) (string, error) {
	data, ok := fetch(data, s.Name)
	if data == nil || !ok {
		return "", nil
	}
	return s.Inner.Render(data)
}

type Comment struct {
	Text string
}

func (c *Comment) Render(data any) (string, error) {
	return "", nil
}

type Reference struct {
	Name   string
	Escape bool
}

var htmlEscaper = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&quot;",
)

func fetch(data any, name string) (any, bool) {
	if name == "." {
		return data, true
	}
	for key := range strings.SplitSeq(name, ".") {
		if data == nil {
			return nil, false
		}
		mapping, ok := data.(map[string]any)
		if !ok {
			return nil, false
		}
		data, ok = mapping[key]
		if !ok {
			return nil, false
		}
	}
	return data, true
}

func (r *Reference) Render(data any) (string, error) {
	data, ok := fetch(data, r.Name)
	value := ""
	if ok && data != nil {
		str, ok := data.(string)
		if ok {
			value = str
		} else {
			value = fmt.Sprint(data)

		}
	}
	if r.Escape {
		value = htmlEscaper.Replace(value)
	}
	return value, nil
}

type Plaintext struct {
	Text string
}

func (p *Plaintext) Render(data any) (string, error) {
	return p.Text, nil
}
