package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"unicode/utf8"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/story"
)

var st *story.Story
var tmpl *template.Template

func init() {
	var err error

	st, err = story.FromJSONFile("gopher.json")
	if err != nil {
		log.Fatalf("received error while parsing story: %s", err)
	}

	tmpl, err = template.ParseFiles("story.html")
	if err != nil {
		log.Fatalf("failed to parse template. Received err: %s", err)
	}
}

// CYOAHandler is an http handler for rendering a CYOA story.
// If the provided url path is "/", the intro chapter will be displayed.
// If the provided url is anything else, the chapter with the given id will be displayed.
// In case of an error (chapter not found, can't render template), status code 300 is returned
func CYOAHandler(w http.ResponseWriter, r *http.Request) {
	var chapter story.Chapter
	var err error

	chapterID := trimFirstRune(r.URL.Path)
	if chapterID != "" {
		chapter, err = st.ChapterByID(chapterID)
	} else {
		chapter, err = st.IntroChapter()
	}

	if err != nil {
		http.Error(
			w, fmt.Sprintf("Failed to find chapter by ID. Received err: %s", err), http.StatusNotFound)

		return
	}

	err = tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Failed to render template: %s", err), http.StatusInternalServerError)

		return
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
