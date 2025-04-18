package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

func sortedNames(parsedFile *ast.File) []string {
	names := make([]string, 0, len(parsedFile.Scope.Objects))
	for object := range parsedFile.Scope.Objects {
		names = append(names, object)
	}
	sort.Strings(names)
	return names
}

func parseFile(fileName string) (*ast.File, error) {
	fs := token.NewFileSet()
	return parser.ParseFile(fs, fileName, nil, 0)
}

func temp() template.FuncMap {
	return template.FuncMap{
		"constant": constant,
		"event":    event,
		"private":  private,
	}
}

var pattern = regexp.MustCompile("([a-z])([A-Z])")

func constCase(name string) string {
	return strings.ToUpper(pattern.ReplaceAllString(name, "${1}_${2}"))
}

func event(name string) bool {
	options := map[string]bool{
		"Connect":    false,
		"Disconnect": false,
		"Event":      false,
		"RateLimit":  false,
		"Interface":  false,
	}
	return !options[name]
}

func constant(name string) string {
	if !event(name) {
		return "__" + constCase(name) + "__"
	}

	return constCase(name)
}

func private(name string) string {
	return strings.ToLower(string(name[0])) + name[1:]
}

func main() {
	parsedFile, err := parseFile("events.go")
	if err != nil {
		log.Fatalf("Error parsing events.go: %s", err)
	}

	names := sortedNames(parsedFile)
	temp := template.Must(template.New("handler").Funcs(temp()).Parse(`...`))

	var buf strings.Builder
	err = temp.Execute(&buf, names)
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}

	src, err := format.Source([]byte(buf.String()))
	if err != nil {
		log.Fatalf("Error formatting source: %s", err)
	}

	dir := filepath.Dir(".")
	err = os.WriteFile(filepath.Join(dir, strings.ToLower("handler.go")), src, 0644)
	if err != nil {
		log.Fatalf("Error writing output: %s", err)
	}
}
