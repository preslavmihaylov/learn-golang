package middleware

import (
	"net/http"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/context"
)

type RequireUser struct{}

func (mw *RequireUser) ApplyFunc(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		usr := context.User(ctx)
		if usr == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		nextHandlerFunc(w, r)
	})
}

func (mw *RequireUser) Apply(nextHandler http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(nextHandler.ServeHTTP)
}
