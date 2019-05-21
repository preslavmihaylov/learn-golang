package main

import (
	"fmt"
	"log"
	"net/http"

	urlshort "github.com/preslavmihaylov/learn-golang/gophercises/ex02-urlshort"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex02-urlshort/config"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex02-urlshort/options"
)

func main() {
	opts := options.ParseArgs()

	redirs := map[string]string{}
	err := config.ReadYAMLConfig(opts.YAMLFilename, redirs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", urlshort.MapHandler(redirs))
}
