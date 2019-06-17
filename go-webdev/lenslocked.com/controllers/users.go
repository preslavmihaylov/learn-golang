package controllers

import (
	"fmt"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	service   *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		service:   us,
	}
}

func (uc *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	err := parseForm(r, &form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	err = uc.service.Create(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "You have signed up successfully")
}

func (uc *Users) Login(w http.ResponseWriter, r *http.Request) {
	lform := LoginForm{}
	err := parseForm(r, &lform)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr, err := uc.service.Authenticate(lform.Email, lform.Password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			http.Error(w, "No such user exists", http.StatusBadRequest)
		case models.ErrWrongPassword:
			http.Error(w, "provided password is incorrect", http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "User(%s %s)", usr.Name, usr.Email)
}
