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
  content string
  requires []string
}

var cases = []RequireTest{
  {`//= require plain`,         []string{"plain"}},
  {`//= require "dquotes"`,     []string{"dquotes"}},
  {`//= require 'squotes'`,     []string{"squotes"}},
  {`//= require trailingsp `,   []string{"trailingsp"}},
  {`//= require extension.ext`, []string{"extension.ext"}},

  {"//= require first\n//= require second", []string{"first", "second"}},
}

func TestRequires(t *testing.T) {
  for _, c := range cases {
    req := findRequires(c.content)
    if !eq(req, c.requires) {
      t.Errorf("edges(%q) = %v, want %v", c.content, req, c.requires)
    }
  }
}
