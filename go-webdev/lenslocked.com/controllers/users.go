package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Users struct {
	NewView *views.View
	service *models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		service: us,
	}
}

func (uc *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := uc.NewView.Render(w, nil)
	if err != nil {
		log.Printf("failed to render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		Name:  form.Name,
		Email: form.Email,
	}

	err = uc.service.Create(&usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, form)
}
