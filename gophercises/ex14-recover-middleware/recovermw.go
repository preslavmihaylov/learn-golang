package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

type RecoverMw struct{}

func (mw *RecoverMw) ApplyFunc(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respWriterMw := &respWriterMw{ResponseWriter: w}
		defer func() {
			if r := recover(); r != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Something went wrong. Please try again...\n")

				if os.Getenv("APP_ENV") == "dev" {
					fmt.Fprintf(w, "\nMessage: %s\n", r)
					fmt.Fprintf(w, "Stack Trace:\n")
					fmt.Fprintf(w, string(debug.Stack()))
				}
			} else {
				respWriterMw.flush()
			}
		}()

		nextHandlerFunc(respWriterMw, r)
	})
}

func (mw *RecoverMw) Apply(nextHandler http.Handler) http.Handler {
	return mw.ApplyFunc(nextHandler.ServeHTTP)
}
