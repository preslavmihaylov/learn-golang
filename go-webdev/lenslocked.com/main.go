package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/config"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/controllers"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/emails"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/middleware"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/rand"
	"golang.org/x/oauth2"
)

func main() {
	var isProductionFlag bool
	flag.BoolVar(&isProductionFlag, "prod", false, "Provide this flag in production. It ensures that "+
		"a config file is provided before the application starts.")
	flag.Parse()

	cfg := config.LoadConfig(isProductionFlag)

	emails.Setup(cfg.Mailgun)
	services, err := models.NewServices(
		models.WithGorm(cfg.Database.Dialect(), cfg.Database.ConnectionInfo()),
		models.WithLogMode(!cfg.Server.IsProduction()),
		models.WithUserService(cfg.Server.Pepper, cfg.Server.HMACKey),
		models.WithGalleryService(),
		models.WithImageService(),
		models.WithOAuthService())
	if err != nil {
		log.Fatalf("failed to create users service: %s", err)
	}
	defer func() {
		err := services.Close()
		if err != nil {
			panic(err)
		}
	}()

	// _ = services.DestructiveReset()
	err = services.AutoMigrate()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	oauthConfigs := map[string]*oauth2.Config{}
	oauthConfigs[models.OAuthDropbox] = &oauth2.Config{
		ClientID:     cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Dropbox.AuthURL,
			TokenURL: cfg.Dropbox.TokenURL,
		},
		RedirectURL: "http://localhost:8080/oauth/dropbox/callback",
	}

	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)
	oauthsC := controllers.NewOAuths(services.OAuth, oauthConfigs)
	staticC := controllers.NewStatic()

	// middlewares
	userMw := middleware.User{UserService: services.User}
	requireUserMw := middleware.RequireUser{}

	randBytes, err := rand.Bytes(32)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %s", err)
	}

	csrfMw := csrf.Protect(randBytes, csrf.Secure(cfg.Server.IsProduction()))

	r.Use(middleware.RequestLogger)
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// oauth routes
	r.HandleFunc("/oauth/{service:[a-z]+}/connect", requireUserMw.ApplyFunc(oauthsC.Connect))
	r.HandleFunc("/oauth/{service:[a-z]+}/callback", requireUserMw.ApplyFunc(oauthsC.Callback))
	r.HandleFunc("/oauth/{service:[a-z]+}/test", requireUserMw.ApplyFunc(oauthsC.DropboxTest))

	// static routes
	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contacts", staticC.ContactsView).Methods("GET")
	r.Handle("/faq", staticC.FAQView).Methods("GET")

	// users routes
	r.HandleFunc("/signup", usersC.GetSignup).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/logout", requireUserMw.ApplyFunc(usersC.Logout)).Methods("POST")
	r.HandleFunc("/signup", usersC.PostSignup).Methods("POST")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.Handle("/forgot_password", usersC.ForgotPasswordView).Methods("GET")
	r.HandleFunc("/forgot_password", usersC.ForgotPassword).Methods("POST")
	r.HandleFunc("/reset_password", usersC.GetResetPassword).Methods("GET")
	r.HandleFunc("/reset_password", usersC.PostResetPassword).Methods("POST")

	// galleries routes
	indexGallery := requireUserMw.ApplyFunc(galleriesC.Index)
	newGallery := requireUserMw.Apply(galleriesC.NewView)
	createGallery := requireUserMw.ApplyFunc(galleriesC.Create)
	editGallery := requireUserMw.ApplyFunc(galleriesC.Edit)
	updateGallery := requireUserMw.ApplyFunc(galleriesC.Update)
	deleteGallery := requireUserMw.ApplyFunc(galleriesC.Delete)
	imageUpload := requireUserMw.ApplyFunc(galleriesC.ImageUpload)
	imageDelete := requireUserMw.ApplyFunc(galleriesC.ImageDelete)

	r.HandleFunc("/galleries", indexGallery).Methods("GET").
		Name(controllers.IndexGalleriesRoute)

	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").
		Name(controllers.ShowGalleryRoute)

	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET").
		Name(controllers.EditGalleryRoute)

	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")

	// images routes
	r.HandleFunc("/galleries/{id:[0-9]+}/images", imageUpload).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", imageDelete).Methods("POST")

	// file servers
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	assetHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetHandler))

	fmt.Printf("Listening on port %d...\n", cfg.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), csrfMw(userMw.Apply(r)))
	if err != nil {
		log.Fatal(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>We could not find the page you were looking for :(</h1>"+
		"<p>please email us if you keep being sent to an invalid page</p>")
}
