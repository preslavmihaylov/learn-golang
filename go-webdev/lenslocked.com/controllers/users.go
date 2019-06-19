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

func (uc *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("remember_token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			fmt.Fprintln(w, "no cookie found")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	usr, err := uc.service.ByRememberToken(c.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	fmt.Fprintln(w, usr.Remember)
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

	err = uc.signIn(w, &usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
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

	err = uc.signIn(w, usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
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

	cookie := http.Cookie{Name: "remember_token", Value: u.Remember, HttpOnly: true}
	http.SetCookie(w, &cookie)

	return nil
}
