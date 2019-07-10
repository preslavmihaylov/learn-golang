package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
)

type View struct {
	Template *template.Template
	layout   string
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

	files = append(layouts, files...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			// this is a stubbed version of the function. It should be provided
			// with an implementation before being invoked.
			return "", errors.New("csrfField is not implemented")
		},
		"pathEscape": func(str string) string {
			return url.PathEscape(str)
		},
	}).ParseFiles(files...)
	if err != nil {
		log.Fatalf("failed to create view: %s", err)
	}

	return &View{
		Template: t,
		layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var viewData Data
	switch data := data.(type) {
	case Data:
		// do nothing
		viewData = data
	default:
		viewData = Data{Yield: data}
	}

	if a := getPersistedAlert(r); a != nil {
		log.Println("found persisted alert. Deleting previous one")
		clearPersistedAlert(w)
		viewData.Alert = a
	}

	ctx := r.Context()
	usr := context.User(ctx)
	viewData.User = usr

	var buf bytes.Buffer

	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(&buf, v.layout, viewData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. If the problem persists, "+
			"please email support@lenslocked.com", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}
