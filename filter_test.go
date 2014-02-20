package monk

import (
	"testing"
)

func TestLessFilter(t *testing.T) {
	before := `.foo {
    .bar {
      width: 100%;
    }
  }`
	after := ".foo .bar{width:100%}\n"
	filter := &LessFilter{}
	if filtered, err := filter.Process(nil, before, "less"); err == nil {
		if filtered != after {
			t.Errorf("LessFilter(%q) = %q, want %q", before, filtered, after)
		}
	} else {
		t.Error(err)
	}
}

func templateFilterCompare(c *Context, t *testing.T, before string, expected string) {
	filter := &TemplateFilter{}
	if filtered, err := filter.Process(c, before, "css"); err == nil {
		if filtered != expected {
			t.Errorf("TemplateFilter(%q) = %q, want %q", before, filtered, expected)
		}
	} else {
		t.Error(err)
	}

}

func TestTemplateFilter(t *testing.T) {
	input := `url('{{url "lolcat.png"}}')`

	fs := NewTestFS()
	context := NewContext(fs)
	context.SearchPath("images")

	fs.File("images/lolcat.png", "LOL!")

	templateFilterCompare(context, t, input, `url('/assets/lolcat.png')`)

	context.Config.AssetRoot = "/a/"
	templateFilterCompare(context, t, input, `url('/a/lolcat.png')`)

	context.Config.Fingerprint = true
	templateFilterCompare(context, t, input,
		`url('/a/lolcat-6cd0dbcbc6ac164f970d9de36ea37634.png')`)
}
