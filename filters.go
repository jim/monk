package monk

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

func filter(content string, extension string) (string, error) {
	switch extension {
	case "bs":
		return strings.Replace(content, "a", "b", -1), nil
	case "fs":
		return strings.Replace(content, "f", "x", -1), nil
	case "coffee":
		return coffeeFilter(content)
	case "less":
		return lessFilter(content)
	case "tmpl":
		return tmplFilter(content)
	}
	return content, nil
}

func coffeeFilter(content string) (string, error) {
	cmd := exec.Command("coffee", "-s", "-c")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func lessFilter(content string) (string, error) {
	cmd := exec.Command("lessc", "-")
	cmd.Stdin = strings.NewReader(content)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func tmplFilter(content string) (string, error) {
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
