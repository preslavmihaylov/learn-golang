package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type View struct {
	tmpl   *template.Template
	layout string
}

var (
	viewsDir      = "views/"
	layoutsDir    = viewsDir + "layouts/"
	tmplExtension = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	layouts, err := filepath.Glob(layoutsDir + "*" + tmplExtension)
	if err != nil {
		log.Fatalf("failed globbing for layouts: %s", err)
	}

	for i := range files {
		files[i] = viewsDir + files[i] + tmplExtension
	}

	files = append(files, layouts...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("failed to create view: %s", err)
	}

	return &View{
		tmpl:   t,
		layout: layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	err := v.tmpl.ExecuteTemplate(w, v.layout, data)
	if err != nil {
		return fmt.Errorf("failed to render view: %s", err)
	}

	return nil
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := v.Render(w, nil)
	if err != nil {
		log.Fatalf("failed to serve http: %s", err)
	}
}
