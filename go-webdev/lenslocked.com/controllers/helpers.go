package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse form: %s", err)
	}

	dec := schema.NewDecoder()

	// using this to ignore the CSRF token key in our forms
	dec.IgnoreUnknownKeys(true)

	err = dec.Decode(dst, r.PostForm)
	if err != nil {
		return fmt.Errorf("failed to decode form: %s", err)
	}

	return nil
}
