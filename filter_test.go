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
