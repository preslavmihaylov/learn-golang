package main

import (
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/cyoa"
)

func main() {
	_, err := cyoa.ParseStory("gopher.json")
	if err != nil {
		log.Fatalf("received error while parsing story: %s", err)
	}
}
