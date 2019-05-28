package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex03-cyoa/cmd/web/handlers"
)

func main() {
	http.HandleFunc("/", handlers.CYOAHandler)

	port := 8080
	fmt.Printf("Server listening on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("failed to serve server at %d. Received err: %s", port, err)
	}
}
