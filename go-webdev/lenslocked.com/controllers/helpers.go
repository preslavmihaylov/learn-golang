package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse form: %s", err)
	}

	return parseValues(r.PostForm, dst)
}

func parseURLParams(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse form: %s", err)
	}

	return parseValues(r.Form, dst)
}

func parseValues(values url.Values, dst interface{}) error {
	dec := schema.NewDecoder()

	// using this to ignore the CSRF token key in our forms
	dec.IgnoreUnknownKeys(true)

	err := dec.Decode(dst, values)
	if err != nil {
		return fmt.Errorf("failed to decode form: %s", err)
	}

	return nil
}
