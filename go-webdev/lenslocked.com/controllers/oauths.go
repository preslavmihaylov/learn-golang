package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	llcontext "github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/dropbox"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
	"golang.org/x/oauth2"
)

type OAuths struct {
	service models.OAuthService
	configs map[string]*oauth2.Config
}

func NewOAuths(os models.OAuthService, configs map[string]*oauth2.Config) *OAuths {
	return &OAuths{
		service: os,
		configs: configs,
	}
}

func (o *OAuths) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 5),
	}
	http.SetCookie(w, &cookie)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuths) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := o.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	receivedState := r.FormValue("state")
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "Missing state", http.StatusBadRequest)
		return
	} else if stateCookie == nil || receivedState != stateCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	oauthCode := r.FormValue("code")
	tok, err := oauthConfig.Exchange(context.Background(), oauthCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr := llcontext.User(r.Context())
	if usr == nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var oauth models.OAuth
	oauth.Token = *tok
	oauth.Service = service
	oauth.UserID = usr.ID

	persistedOAuth, err := o.service.Find(oauth.UserID, oauth.Service)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			// good. no previous token exists
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if persistedOAuth != nil {
		err = o.service.Delete(persistedOAuth.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = o.service.Create(&oauth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (o *OAuths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usr := llcontext.User(r.Context())
	if usr == nil {
		panic("endpoint should have require user middleware setup")
	}

	dbxPath := r.FormValue("path")
	usrOAuth, err := o.service.Find(usr.ID, "dropbox")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folders, files, err := dropbox.List(usrOAuth.AccessToken, r.FormValue(dbxPath))
	fmt.Fprintln(w, "Folders=", folders)
	fmt.Fprintln(w, "Files=", files)
}
