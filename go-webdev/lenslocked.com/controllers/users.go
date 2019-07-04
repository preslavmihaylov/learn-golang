package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Users struct {
	SignupView *views.View
	LoginView  *views.View
	service    models.UserService
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

func NewUsers(us models.UserService) *Users {
	return &Users{
		SignupView: views.NewView("bootstrap", "users/signup"),
		LoginView:  views.NewView("bootstrap", "users/login"),
		service:    us,
	}
}

func (uc *Users) Create(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data

	form := SignupForm{}
	err := parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		uc.SignupView.Render(w, r, viewData)
		return
	}

	usr := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	err = uc.service.Create(&usr)
	if err != nil {
		viewData.SetAlert(err)
		uc.SignupView.Render(w, r, viewData)
		return
	}

	err = uc.signIn(w, &usr)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (uc *Users) Login(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data

	lform := LoginForm{}
	err := parseForm(r, &lform)
	if err != nil {
		viewData.SetAlert(err)
		uc.LoginView.Render(w, r, viewData)
		return
	}

	usr, err := uc.service.Authenticate(lform.Email, lform.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			viewData.AlertError("No user exists with that email address")
		default:
			viewData.SetAlert(err)
		}

		uc.LoginView.Render(w, r, viewData)
		return
	}

	err = uc.signIn(w, usr)
	if err != nil {
		viewData.SetAlert(err)
		uc.LoginView.Render(w, r, viewData)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (uc *Users) Logout(w http.ResponseWriter, r *http.Request) {
	usr := context.User(r.Context())
	uc.logout(w, usr)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (uc *Users) signIn(w http.ResponseWriter, u *models.User) error {
	err := u.GenerateToken()
	if err != nil {
		return err
	}

	err = uc.service.Update(u)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	cookie := http.Cookie{Name: "remember_token", Value: u.RememberToken, HttpOnly: true}
	http.SetCookie(w, &cookie)

	return nil
}

func (uc *Users) logout(w http.ResponseWriter, u *models.User) error {
	err := u.GenerateToken()
	if err != nil {
		return err
	}

	err = uc.service.Update(u)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	cookie := http.Cookie{Name: "remember_token", Value: "", Expires: time.Now(), HttpOnly: true}
	http.SetCookie(w, &cookie)

	return nil
}
