package middleware

import (
	"net/http"
	"strings"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
)

type User struct {
	models.UserService
}

func (mw *User) ApplyFunc(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			nextHandlerFunc(w, r)
			return
		}

		c, err := r.Cookie("remember_token")
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				nextHandlerFunc(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		usr, err := mw.UserService.ByRememberToken(c.Value)
		if err != nil {
			switch err {
			case models.ErrNotFound:
				nextHandlerFunc(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// store the logged in user in the request context
		ctx := r.Context()
		ctx = context.WithUser(ctx, usr)
		r = r.WithContext(ctx)

		nextHandlerFunc(w, r)
	})
}

func (mw *User) Apply(nextHandler http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(nextHandler.ServeHTTP)
}
