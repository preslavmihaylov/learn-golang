package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/emails"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/views"
)

type Users struct {
	SignupView         *views.View
	LoginView          *views.View
	ForgotPasswordView *views.View
	ResetPasswordView  *views.View
	service            models.UserService
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

type ResetPasswordForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService) *Users {
	return &Users{
		SignupView:         views.NewView("bootstrap", "users/signup"),
		LoginView:          views.NewView("bootstrap", "users/login"),
		ForgotPasswordView: views.NewView("bootstrap", "users/forgot_password"),
		ResetPasswordView:  views.NewView("bootstrap", "users/reset_password"),
		service:            us,
	}
}

func (uc *Users) GetSignup(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	err := parseURLParams(r, &form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc.SignupView.Render(w, r, form)
}

func (uc *Users) PostSignup(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data
	form := SignupForm{}
	viewData.Yield = &form

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
	err := uc.logout(w, usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (uc *Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data
	var form ResetPasswordForm

	viewData.Yield = &form
	err := parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		uc.ForgotPasswordView.Render(w, r, viewData)
		return
	}

	token, err := uc.service.InitiatePasswordReset(form.Email)
	if err != nil {
		viewData.SetAlert(err)
		uc.ForgotPasswordView.Render(w, r, viewData)
		return
	}

	err = emails.SendResetPasswordEmail(form.Email, token)
	if err != nil {
		viewData.SetAlert(err)
		uc.ForgotPasswordView.Render(w, r, viewData)
		return
	}

	views.RedirectWithAlert(w, r, "/reset_password", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})
}

func (uc *Users) GetResetPassword(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data
	var form ResetPasswordForm

	viewData.Yield = &form
	err := parseURLParams(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		uc.ResetPasswordView.Render(w, r, viewData)
		return
	}

	uc.ResetPasswordView.Render(w, r, viewData)
}

func (uc *Users) PostResetPassword(w http.ResponseWriter, r *http.Request) {
	var viewData views.Data
	var form ResetPasswordForm

	err := parseForm(r, &form)
	if err != nil {
		viewData.SetAlert(err)
		uc.ResetPasswordView.Render(w, r, viewData)
		return
	}

	usr, err := uc.service.CompletePasswordReset(form.Token, form.Password)
	if err != nil {
		viewData.SetAlert(err)
		uc.ResetPasswordView.Render(w, r, viewData)
		return
	}

	err = uc.signIn(w, usr)
	if err != nil {
		viewData.SetAlert(err)
		uc.ResetPasswordView.Render(w, r, viewData)
		return
	}

	views.RedirectWithAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "You have successfully reset your password!",
	})
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
