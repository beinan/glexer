package main

import (
	"html/template"
	"os"
)

func main() {
	d := struct {
		Name string
	}{
		"aa",
	}

	t := template.Must(template.New("queue").Parse(rootResolverTemplate))
	t.Execute(os.Stdout, d)
}

var rootResolverTemplate = `
package graphql

{{.Name}}

`
