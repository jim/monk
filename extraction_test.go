package monk

import (
	"testing"
)

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type RequireTest struct {
	rawContent string
	content    string
	requires   []string
}

var cases = []RequireTest{
	{`//= require plain`, "", []string{"plain"}},
	{`//= require "dquotes"`, "", []string{"dquotes"}},
	{`//= require 'squotes'`, "", []string{"squotes"}},
	{`//= require trailingsp `, "", []string{"trailingsp"}},
	{`//= require extension.ext`, "", []string{"extension.ext"}},
	{`//=require nospace`, "", []string{"nospace"}},
	{`//=  require  manyspace`, "", []string{"manyspace"}},

	{"//= require first\n//= require second", "", []string{"first", "second"}},

	{"//= require dep\nother content", "other content", []string{"dep"}},
}

func TestRequires(t *testing.T) {
	for _, c := range cases {
		content, requires := extractDependencies(c.rawContent)
		if !eq(requires, c.requires) {
			t.Errorf("extractDependencies(%q) requires = %v, want %v", c.content, requires, c.requires)
		}
		if content != c.content {
			t.Errorf("extractDependencies(%q) content = %q, want %q", c.rawContent, content, c.content)
		}
	}
}
