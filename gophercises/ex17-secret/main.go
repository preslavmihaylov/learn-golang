package main

import (
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex17-secret/secret"
)

func main() {
	v, err := secret.FileVault("encoding-key", "secrets")
	if err != nil {
		log.Fatalf("failed to create file vault: %s", err)
	}

	v.Set("key-name", "key-value")
	v.Set("key-2", "key-2")

	value, err := v.Get("key-name")
	fmt.Println(value) // "key-value"

	value, err = v.Get("key-2")
	fmt.Println(value) // "key-value"
}
