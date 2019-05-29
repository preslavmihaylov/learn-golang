package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/cmd/web/cyoa_handlers"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/options"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/story"
)

func main() {
	defaultStory, err := story.FromJSONFile(options.JSONFilename())
	if err != nil {
		log.Fatalf("received error while parsing default story: %s", err)
	}

	defaultTmpl, err := template.ParseFiles(options.HTMLTemplateFilename())
	if err != nil {
		log.Fatalf("failed to parse default template. Received err: %s", err)
	}

	cyoaHandler := cyoa_handlers.NewCYOAHandler(defaultStory, defaultTmpl)

	port := 8080
	fmt.Printf("Server listening on port %d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), cyoaHandler)
	if err != nil {
		log.Fatalf("failed to serve server at %d. Received err: %s", port, err)
	}
}
