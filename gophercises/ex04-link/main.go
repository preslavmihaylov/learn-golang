package main

import (
	"fmt"
	"log"
	"os"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex04-link/html/links"
)

func main() {
	test1()
	test2()
	test3()
	test4()
}

func test1() {
	fmt.Println("------ TEST 1 ------")
	test("ex1.html")
	fmt.Println()
}

func test2() {
	fmt.Println("------ TEST 2 ------")
	test("ex2.html")
	fmt.Println()
}

func test3() {
	fmt.Println("------ TEST 3 ------")
	test("ex3.html")
	fmt.Println()
}

func test4() {
	fmt.Println("------ TEST 4 ------")
	test("ex4.html")
	fmt.Println()
}

func test(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open html file: %s", err)
	}

	lnks, err := links.Parse(f)
	if err != nil {
		log.Fatalf("received error while parsing links: %s", err)
	}

	for _, lnk := range lnks {
		fmt.Printf("%s -> %s\n", lnk.Href, lnk.Text)
	}
}
