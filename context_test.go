package monk

import (
	"strings"
	"testing"
)

func TestFindAssetInSearchPaths(t *testing.T) {
	fs := NewTestFS()
	ac := NewContext(fs)

	fs.File("assets/simple.js", "")
	needle := "simple"

	_, err := ac.findAssetInSearchPaths(needle)

	if err == nil || !strings.Contains(err.Error(), "No search paths") {
		t.Errorf("should have required at least one search path to be defined, got: %s", err)
	}

	ac.SearchPath("assets")
	_, err = ac.findAssetInSearchPaths(needle)

	if err == nil || !strings.Contains(err.Error(), "extension is required") {
		t.Errorf("should have required %s to have an extension", needle)
	}
}

func TestCreateAssetContent(t *testing.T) {
	fs := NewTestFS()
	c := NewContext(fs)

	assetContent := "//= require a"
	assetPath := "assets/simple.js"

	fs.File(assetPath, assetContent)
	c.SearchPath("assets")
	expected := []string{"a.js"}

	if info, err := fs.Stat(assetPath); err != nil {
		t.Errorf("Unable to stat %q: %s", assetPath, err)
	} else {

		if asset, err := c.createAsset(assetPath, info); err != nil {
			t.Errorf("Tried to create an asset, got: %s", err)
		} else {
			if !eq(asset.Dependencies, expected) {
				t.Errorf("Expected dependencies to be %q, got %q", expected, asset.Dependencies)
			}
		}

	}
}

func TestLoadAssetContent(t *testing.T) {
	fs := NewTestFS()
	ac := NewContext(fs)
	assetContent := "//= require a"
	assetPath := "assets/simple.js"

	fs.File(assetPath, assetContent)
	ac.SearchPath("assets")

	if content, err := ac.loadAssetContent(assetPath); err == nil {
		if content != assetContent {
			t.Errorf("requiring %q, want %q, got %q", assetPath, assetContent, content)
		}
	} else {
		t.Error(err)
	}
}

func TestSearchDirectory(t *testing.T) {
	fs := NewTestFS()
	ac := NewContext(fs)
	assetPath := "assets/simple.js"
	query := "simple"

	fs.File(assetPath, "")
	ac.SearchPath("assets")

	if absPath, err := ac.searchDirectory("assets", query); err == nil {
		if absPath != assetPath {
			t.Errorf("searchDirectory(%q), want %q, got %q", query, assetPath, absPath)
		}
	} else {
		t.Error(err)
	}
}

func TestExplodeDepencies(t *testing.T) {
	fs := NewTestFS()
	fs.File("bar/1", "")
	fs.File("bar/2", "")
	fs.File("baz", "")

	d := []string{"foo", "bar/*", "baz"}
	dirPath := "assets"
	exploded := explodeDependencies(dirPath, d, fs)
	expected := []string{"foo", "bar/1", "bar/2", "baz"}
	if !eq(expected, exploded) {
		t.Errorf("explodeDependencies(%v) = %v, want %v", d, exploded, expected)
	}
}
