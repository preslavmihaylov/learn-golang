package cyoa_handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/story"
)

type cyoaHandler struct {
	story  *story.Story
	tmpl   *template.Template
	prefix string
}

// CYOAHandlerOpts is a functional option wrapper for the internal cyoaHandler type
type CYOAHandlerOpts func(h *cyoaHandler)

// NewCYOAHandler constructs a new cyoa handler implementing the http.Handler interface.
func NewCYOAHandler(st *story.Story, tmpl *template.Template, opts ...CYOAHandlerOpts) http.Handler {
	h := cyoaHandler{st, tmpl, "/"}
	for _, opt := range opts {
		opt(&h)
	}

	return &h
}

// WithPrefix returns a functional option setting the prefix of the internal cyoa handler type
func WithPrefix(prefix string) CYOAHandlerOpts {
	return func(h *cyoaHandler) {
		h.prefix = prefix
	}
}

// ServeHTTP renders a CYOA story.
// If the provided url path is the handler's prefix (default "/"), the intro chapter will be displayed.
// If the provided url is anything else, the chapter with the given id will be displayed.
// In case of an error (chapter not found, can't render template), an error status code is returned
func (h *cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var chapter story.Chapter
	var err error

	chapterID := trimPrefix(r.URL.Path, h.prefix)
	if chapterID != "" {
		chapter, err = h.story.ChapterByID(chapterID)
	} else {
		chapter, err = h.story.IntroChapter()
	}

	if err != nil {
		http.Error(
			w, fmt.Sprintf("Failed to find chapter by ID. Received err: %s", err), http.StatusNotFound)

		return
	}

	err = h.tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)

		return
	}
}

func trimPrefix(src string, prefix string) string {
	return src[len(prefix):]
}
