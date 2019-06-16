package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/controllers"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "developer"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	usersS, err := models.NewUserService(psqlInfo)
	if err != nil {
		log.Fatalf("failed to create users service: %s", err)
	}
	defer func() {
		err := usersS.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = usersS.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	usersC := controllers.NewUsers(usersS)
	staticC := controllers.NewStatic()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contacts", staticC.ContactsView).Methods("GET")
	r.Handle("/faq", staticC.FAQView).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

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
