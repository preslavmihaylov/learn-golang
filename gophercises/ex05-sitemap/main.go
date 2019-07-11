package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex05-sitemap/html/sitemap"
)

func main() {
	domain := flag.String("domain", "", "domain to start from")
	outputXML := flag.String("xml", "", "output xml")
	flag.Parse()

	if *domain == "" {
		fmt.Println("Invalid domain name provided. See help for options.")
		return
	} else if *outputXML == "" {
		fmt.Println("Invalid xml file provided. See help for options")
		return
	}

	f, err := os.Create(*outputXML)
	if err != nil {
		log.Fatalf("failed to open output xml file: %s", err)
	}

	err = sitemap.Build(f, *domain)
	if err != nil {
		log.Fatalf("failed building sitemap: %s", err)
	}
}
