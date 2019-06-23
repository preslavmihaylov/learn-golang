package views

import (
	"bytes"
	"html/template"
	"io"
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

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{Yield: data}
	}

	var buf bytes.Buffer
	err := v.tmpl.ExecuteTemplate(&buf, v.layout, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, "+
			"please email support@lenslocked.com", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}
