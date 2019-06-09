package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)

	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to my awesome site</h1>")
	} else if r.URL.Path == "/contacts" {
		fmt.Fprintf(w, "to get in touch, please send an email to "+
			"<a href=\"mailto:support@lenslocked.com\">support@lenslocked.com</a>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>We could not find the page you were looking for :(</h1>"+
			"<p>please email us if you keep being sent to an invalid page</p>")
	}
}
