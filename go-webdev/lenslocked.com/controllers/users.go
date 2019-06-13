package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers() *Users {
	return &Users{NewView: views.NewView("bootstrap", "users/new")}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := u.NewView.Render(w, nil)
	if err != nil {
		log.Printf("failed to render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignupForm{}
	err := parseForm(r, &form)
	if err != nil {
		log.Printf("failed to parse form: %s", err)
		fmt.Fprintf(w, "Invalid Form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, form)
}
