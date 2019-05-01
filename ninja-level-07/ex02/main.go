package main

import "fmt"

type person struct {
	name string
}

func changeMe(p *person) {
	p.name = "George"
}

func main() {
	p := person{"Peter"}

	fmt.Println(p.name)
	changeMe(&p)
	fmt.Println(p.name)
}
