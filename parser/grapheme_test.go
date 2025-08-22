package parser

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *Grapheme
	}{
		{"Normal New", args{"abc"}, &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}},
		{"Empty New", args{""}, &Grapheme{"", "", 1, 0, 0, -1, 0}},
		{"New Line New", args{"\nabc"}, &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGrapheme(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphemeNext(t *testing.T) {
	tests := []struct {
		name string
		g    *Grapheme
		want *Grapheme
	}{
		{"Initial Next", &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}, &Grapheme{"b", "c", 1, 2, 2, 2320992, 16}},
		{"Normal Next", &Grapheme{"a", "bc", 3, 17, 41, 2320992, 16}, &Grapheme{"b", "c", 3, 18, 42, 2320992, 16}},
		{"End Next", &Grapheme{"a", "", 7, 1, 53, 2320992, 16}, &Grapheme{"", "", 7, 2, 54, -1, 0}},
		{"After End Next", &Grapheme{"", "", 5, 8, 14, -1, 0}, &Grapheme{"", "", 5, 8, 14, -1, 0}},
		{"New Line Next", &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}, &Grapheme{"a", "bc", 2, 1, 2, 2320992, 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Next(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grapheme.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphemeIsEof(t *testing.T) {
	tests := []struct {
		name string
		g    *Grapheme
		want bool
	}{
		{"Initial IsEof", &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}, false},
		{"Normal IsEof", &Grapheme{"a", "bc", 3, 17, 41, 2320992, 16}, false},
		{"End IsEof", &Grapheme{"a", "", 7, 1, 53, 2320992, 16}, false},
		{"After End IsEof", &Grapheme{"", "", 5, 8, 14, -1, 0}, true},
		{"New Line IsEof", &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.IsEof(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grapheme.IsEof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphemeIsEol(t *testing.T) {
	tests := []struct {
		name string
		g    *Grapheme
		want bool
	}{
		{"Initial IsEol", &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}, false},
		{"Normal IsEol", &Grapheme{"a", "bc", 3, 17, 41, 2320992, 16}, false},
		{"End IsEol", &Grapheme{"a", "", 7, 1, 53, 2320992, 16}, false},
		{"After End IsEol", &Grapheme{"", "", 5, 8, 14, -1, 0}, false},
		{"New Line IsEol", &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.IsEol(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grapheme.IsEol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrapheme_String(t *testing.T) {
	tests := []struct {
		name string
		g    *Grapheme
		want string
	}{
		{"Initial String", &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}, "'a' 1:1 (1)"},
		{"Normal String", &Grapheme{"a", "bc", 3, 17, 41, 2320992, 16}, "'a' 3:17 (41)"},
		{"End String", &Grapheme{"a", "", 7, 1, 53, 2320992, 16}, "'a' 7:1 (53)"},
		{"After End String", &Grapheme{"", "", 5, 8, 14, -1, 0}, "EOF 5:8 (14)"},
		{"New Line String", &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}, "'\n' 1:1 (1)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("Grapheme.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrapheme_Error(t *testing.T) {
	type args struct {
		expected string
	}
	tests := []struct {
		name   string
		g      *Grapheme
		args   args
		errMsg string
	}{
		{"Normal Error", &Grapheme{"a", "bc", 3, 2, 15, 2320992, 16}, args{"b"}, "at 'a' 3:2 (15) expected b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.g.Error(tt.args.expected); err.Error() != tt.errMsg {
				t.Errorf("Grapheme.Error() error = %v, wantErr %v", err, tt.errMsg)
			}
		})
	}
}

func TestGraphemes(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []*Grapheme
	}{
		{"Normal Graphemes", args{"a\nbc"}, []*Grapheme{{"a", "\nbc", 1, 1, 1, 8414242, 20}, {"\n", "bc", 1, 2, 2, 2320992, 14}, {"b", "c", 2, 1, 3, 2320992, 16}, {"c", "", 2, 2, 4, 0, 30}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := 0
			for got := range Graphemes(tt.args.str) {
				if !reflect.DeepEqual(got, tt.want[i]) {
					t.Errorf("Graphemes() = %v, want %v", got, tt.want[i])
				}
				i++
			}
		})
	}
}

func TestGraphemesBreak(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []*Grapheme
	}{
		{"Normal Graphemes", args{"a\nbc"}, []*Grapheme{{"a", "\nbc", 1, 1, 1, 8414242, 20}, {"\n", "bc", 1, 2, 2, 2320992, 14}, {"b", "c", 2, 1, 3, 2320992, 16}, {"c", "", 2, 2, 4, 0, 30}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := 0
			for got := range Graphemes(tt.args.str) {
				if !reflect.DeepEqual(got, tt.want[i]) {
					t.Errorf("Graphemes() = %v, want %v", got, tt.want[i])
				}
				i++
				if i == 2 {
					break
				}
			}
		})
	}
}
