package monk

import (
	"regexp"
)

// Parse an asset's dependencies and return the content stripped of these declarations,
// along with a slice containing the dependencies that were declared.
func extractDependencies(fileContents string) (string, []string) {
	pattern := `(?m)^\s*//=\s*require\s+['"]?([\w\.]+)["']?\s*$?`
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	requires := make([]string, 0)
	stripped := r.ReplaceAllStringFunc(fileContents, func(line string) string {
		match := r.FindStringSubmatch(line)
		requires = append(requires, match[1])
		return ""
	})

	return stripped, requires
}
