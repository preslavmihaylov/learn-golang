package main

import "fmt"

func main() {
	person := struct {
		firstName string
		lastName  string
		age       int
	}{
		firstName: "Page",
		lastName:  "Two",
		age:       2,
	}

	fmt.Println(person)
}
