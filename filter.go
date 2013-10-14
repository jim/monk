package monk

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

var filters = map[string]AssetProcessor{}

type AssetProcessor interface {
	Process(content string, extension string) (string, error)
	CheckSystem() error
}

type AssetFilter struct{}

func (af AssetFilter) CheckSystem() error {
	return nil
}

func (af AssetFilter) RequireBin(bin string) error {
	cmd := exec.Command("which", bin)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil || strings.Contains(out.String(), "not found") {
		return fmt.Errorf("The command %q was not found on your path", bin)
	}
	return nil
}

func init() {
	AppendFilter("coffee", &CoffeeFilter{})
	AppendFilter("less", &LessFilter{})
	AppendFilter("tmpl", &TemplateFilter{})
}

func AppendFilter(extension string, filter AssetProcessor) {
	if err := filter.CheckSystem(); err != nil {
		panic(err)
	}
	filters[extension] = filter
}

type CoffeeFilter struct {
	AssetFilter
}

func (cf CoffeeFilter) Process(content string, extension string) (string, error) {
	cmd := exec.Command("coffee", "-s", "-c")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func (cf CoffeeFilter) CheckSystem() error {
	return cf.RequireBin("coffee")
}

type LessFilter struct {
	AssetFilter
}

func (lf LessFilter) Process(content string, extension string) (string, error) {
	cmd := exec.Command("lessc", "-", "--compress")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func (a AssetProcessor) Foo() string {
	return "foo"
}

func (lf LessFilter) CheckSystem() error {
	return lf.RequireBin("lessc")
}

type TemplateFilter struct{}

func (tf TemplateFilter) Process(content string, extension string) (string, error) {
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

func (tf TemplateFilter) CheckSystem() error {
	return nil
}

func ApplyFilter(content string, extension string) (string, error) {
	if filter, ok := filters[extension]; ok {
		return filter.Process(content, extension)
	}
	return "", fmt.Errorf("could not find a filter for extension: %q", extension)
}
