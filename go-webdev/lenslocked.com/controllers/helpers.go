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
	err = dec.Decode(dst, r.PostForm)
	if err != nil {
		return fmt.Errorf("failed to decode form: %s", err)
	}

	return nil
}
