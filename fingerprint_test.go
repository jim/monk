package monk

import (
	"testing"
)

func TestGenerateFingerprint(t *testing.T) {

	fs := NewTestFS()
	fs.File("notes", "all in due time\n")

	fingerprint, err := GenerateFingerprint(fs, "notes")

	if err != nil {
		t.Fatal(err)
	}

	expected := "858cca6356a811039b2367bdbd2acaee"
	if fingerprint != expected {
		t.Errorf("expected %q, got: %q", expected, fingerprint)
	}
}
