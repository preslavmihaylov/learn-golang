package main

import "fmt"

type person struct {
	firstName string
	lastName  string
	age       int
}

func (p person) speak() {
	fmt.Println("Name:", p.firstName, p.lastName, ", Age:", p.age)
}

func main() {
	p := person{
		firstName: "Peter",
		lastName:  "Jackson",
		age:       30,
	}

	p.speak()
}
