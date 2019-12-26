package main

import (
	"log"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex15-recover-chroma/handlers"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex15-recover-chroma/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", handlers.PanicHandler)
	mux.HandleFunc("/panic-after/", handlers.PanicAfterHandler)
	mux.HandleFunc("/debug", handlers.DebugHandler)
	mux.HandleFunc("/", handlers.HomeHandler)
	log.Fatal(http.ListenAndServe(":3000", middleware.Development(mux)))
}
