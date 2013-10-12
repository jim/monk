package monk

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

type AssetFilter func (content string, extension string) (string, error)

var filters = map[string]AssetFilter{}

func init() {
  AppendFilter("coffee", CoffeeFilter)
  AppendFilter("less", LessFilter)
  AppendFilter("tmpl", TemplateFilter)
}

func AppendFilter(extension string, filter AssetFilter) {
  filters[extension] = filter
}

func CoffeeFilter(content string, extension string) (string, error) {
	cmd := exec.Command("coffee", "-s", "-c")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func LessFilter(content string, extension string) (string, error) {
	cmd := exec.Command("lessc", "-", "--compress")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func TemplateFilter(content string, extension string) (string, error) {
	tmpl := template.New("asset")

	helpers := template.FuncMap{
		"url": func(assetPath string) (string, error) {
			return fmt.Sprintf("%s-%s", assetPath, "fingerprint"), nil
		},
	}

	tmpl.Funcs(helpers)
	_, err := tmpl.Parse(content)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, nil)

	return out.String(), err
}

func ApplyFilters(content string, extension string) (string, error) {
  for ext, filter := range filters {
    if ext == extension {
      return filter(content, extension)
    }
  }
  return "", fmt.Errorf("could not find a filter for extension: %q", extension)
}
