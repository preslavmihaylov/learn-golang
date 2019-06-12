package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

var homeView *views.View
var contactsView *views.View
var faqView *views.View

func main() {
	var err error
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactsView = views.NewView("bootstrap", "views/contacts.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/contacts", contactsHandler)
	r.HandleFunc("/faq", faqHandler)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	fmt.Println("Listening on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>We could not find the page you were looking for :(</h1>"+
		"<p>please email us if you keep being sent to an invalid page</p>")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Render(w, nil)
	if err != nil {
		log.Printf("home handler failed to render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactsView.Render(w, nil)
	if err != nil {
		log.Printf("contacts handler failed to render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := faqView.Render(w, nil)
	if err != nil {
		log.Printf("faq handler failed to render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
