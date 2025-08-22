package parser

import (
	"fmt"

	"github.com/rivo/uniseg"
)

type Grapheme struct {
	Token, remaining                     string
	Line, Column, Pos, state, boundaries int
}

func NewGrapheme(str string) *Grapheme {
	return (&Grapheme{"", str, 1, 0, 0, -1, 0}).Next()
}

func (g *Grapheme) Next() *Grapheme {
	if g.remaining == "" {
		if g.Token == "" {
			return g
		}
		return &Grapheme{"", "", g.Line, g.Column + 1, g.Pos + 1, -1, 0}
	}
	ch, remaining, boundaries, state := uniseg.StepString(g.remaining, g.state)
	if g.IsEol() {
		return &Grapheme{ch, remaining, g.Line + 1, 1, g.Pos + 1, state, boundaries}
	}
	return &Grapheme{ch, remaining, g.Line, g.Column + 1, g.Pos + 1, state, boundaries}
}

func (g *Grapheme) IsEof() bool {
	return g.Token == ""
}

func (g *Grapheme) IsEol() bool {
	return g.boundaries&uniseg.MaskLine == uniseg.LineMustBreak
}

func (g *Grapheme) String() string {
	if g.IsEof() {
		return fmt.Sprintf("EOF %d:%d (%d)", g.Line, g.Column, g.Pos)
	}
	return fmt.Sprintf("'%s' %d:%d (%d)", g.Token, g.Line, g.Column, g.Pos)
}

func (g *Grapheme) Error(expected string) error {
	return fmt.Errorf("at %s expected %s", g, expected)
}

func Graphemes(str string) func(func(*Grapheme) bool) {
	return func(yield func(*Grapheme) bool) {
		for seq := NewGrapheme(str); !seq.IsEof(); seq = seq.Next() {
			if !yield(seq) {
				return
			}
		}
	}
}
