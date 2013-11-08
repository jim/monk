package monk

import (
	"testing"
)

var expected = `/* d.js */
source of d

/* f.js */
source of f

/* e.js */
doSomething('path-to-asset.jpg-fingerprint');

/* c.js */

/* b.js */
source of b

/* g.js */
// Generated by CoffeeScript 1.3.3
(function() {
  var Foo;

  Foo = (function() {

    function Foo() {}

    Foo.prototype.bar = function(i) {
      return console.log(i);
    };

    return Foo;

  })();

}).call(this);

/* a.js */

`

func TestBuild(t *testing.T) {

	fs := NewTestFS()

	fs.File("assets/a.js", "//= require b\n//= require d\n//= require g\n")
	fs.File("assets/b.js", "//= require c\n//= require c\n\nsource of b\n")
	fs.File("assets/c.js", "//= require d\n//= require e\n")
	fs.File("assets/d.js", "source of d\n")
	fs.File("assets/e.js.tmpl", "//= require f\n\n"+`doSomething('{{url "path-to-asset.jpg" }}');`+"\n")
	fs.File("assets/f.js", "source of f\n")
	fs.File("assets/g.js.coffee", "class Foo\n  bar: (i) -> console.log i\n")

	context := NewContext(fs)
	context.SearchPath("assets")

	r := &Resolution{}

	if err := r.Resolve("a.js", context); err != nil {
		t.Fatal(err)
	}

	built := Build(r, context)
	if built != expected {
		t.Errorf("expected %q, got: %q", expected, built)
	}
}