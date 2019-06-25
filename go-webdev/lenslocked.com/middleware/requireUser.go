package middleware

import (
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplyFunc(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("remember_token")
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				http.Redirect(w, r, "/login", http.StatusFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		usr, err := mw.UserService.ByRememberToken(c.Value)
		if err != nil {
			switch err {
			case models.ErrNotFound:
				http.Redirect(w, r, "/login", http.StatusFound)
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

func (mw *RequireUser) Apply(nextHandler http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(nextHandler.ServeHTTP)
}
