package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/controllers"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/middleware"
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

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		log.Fatalf("failed to create users service: %s", err)
	}
	defer func() {
		err := services.Close()
		if err != nil {
			panic(err)
		}
	}()

	// _ = usersS.DestructiveReset()
	err = services.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Images, r)
	staticC := controllers.NewStatic()

	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}

	r.Use(middleware.RequestLogger)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// static routes
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contacts", staticC.ContactsView).Methods("GET")
	r.Handle("/faq", staticC.FAQView).Methods("GET")

	// users routes
	r.Handle("/signup", usersC.SignupView).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	// galleries routes
	indexGallery := requireUserMw.ApplyFunc(galleriesC.Index)
	newGallery := requireUserMw.Apply(galleriesC.NewView)
	createGallery := requireUserMw.ApplyFunc(galleriesC.Create)
	showGallery := requireUserMw.ApplyFunc(galleriesC.Show)
	editGallery := requireUserMw.ApplyFunc(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFunc(galleriesC.Update)
	deleteGallery := requireUserMw.ApplyFunc(galleriesC.Delete)
	imageUpload := requireUserMw.ApplyFunc(galleriesC.ImageUpload)
	imageDelete := requireUserMw.ApplyFunc(galleriesC.ImageDelete)

	r.HandleFunc("/galleries", indexGallery).Methods("GET").
		Name(controllers.IndexGalleriesRoute)

	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", showGallery).Methods("GET").
		Name(controllers.ShowGalleryRoute)

	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET").
		Name(controllers.EditGalleryRoute)

	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")

	// images routes
	r.HandleFunc("/galleries/{id:[0-9]+}/images", imageUpload).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", imageDelete).Methods("POST")

	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	fmt.Println("Listening on port 8080...")
	err = http.ListenAndServe(":8080", userMw.Apply(r))
	if err != nil {
		log.Fatal(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>We could not find the page you were looking for :(</h1>"+
		"<p>please email us if you keep being sent to an invalid page</p>")
}
