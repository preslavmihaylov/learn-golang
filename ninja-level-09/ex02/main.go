package main

import "fmt"

type human interface {
	speak()
}

type person struct {
	firstName string
	lastName  string
}

func (p *person) speak() {
	fmt.Println("I am", p.firstName, p.lastName)
}

func saySomething(h human) {
	h.speak()
}

func main() {
	p := person{firstName: "George", lastName: "Angelov"}
	saySomething(&p)

	// does not compile
	// saySomething(p)
}
