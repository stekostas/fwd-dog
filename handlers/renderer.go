package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Renderer struct {
	StaticFilesRoot string
	AssetsRoot      string
}

func NewRenderer(staticFilesRoot string, assetsRoot string) *Renderer {
	return &Renderer{StaticFilesRoot: staticFilesRoot, AssetsRoot: assetsRoot}
}

func (r *Renderer) RenderTemplate(name string, data interface{}, writer http.ResponseWriter) {
	templateName := fmt.Sprintf("%s/%s", r.StaticFilesRoot, name)
	tmpl := template.Must(template.ParseFiles(templateName))
	err := tmpl.Execute(writer, data)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Could not serve template '%s': %v\n", templateName, err)
	}
}
