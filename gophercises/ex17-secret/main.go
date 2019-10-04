package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex17-secret/secret"
)

func main() {
	currUsr, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get current user: %s", err)
	}

	cmd := flag.String("c", "", "command you want to execute. [set, get]")
	key := flag.String("k", "", "key you want to set")
	val := flag.String("v", "", "value you want to set")
	pass := flag.String("p", "", "the passphrase of file")
	inputFilename := flag.String("i", currUsr.HomeDir+"/.secrets_goph", "input file. Defaults to ~/.secrets_goph if not specified")
	flag.Parse()

	if *cmd == "" {
		log.Fatalf("command is required")
	} else if *key == "" {
		log.Fatalf("key is required")
	} else if *cmd == "set" && *val == "" {
		log.Fatalf("value is required")
	} else if *pass == "" {
		log.Fatalf("passphrase is required")
	}

	v, err := secret.FileVault(*pass, *inputFilename)
	if err != nil {
		log.Fatalf("failed to create file vault: %s", err)
	}

	switch *cmd {
	case "set":
		err = v.Set(*key, *val)
		if err != nil {
			log.Fatalf("failed to set key: %s", err)
		}

		fmt.Println("the key has been saved")
	case "get":
		value, err := v.Get(*key)
		if err != nil {
			log.Fatalf("failed to get key: %s", err)
		}

		fmt.Println(value)
	}
}
