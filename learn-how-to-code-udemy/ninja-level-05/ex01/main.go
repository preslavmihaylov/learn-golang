package main

import "fmt"

type iceCreamFlavor int

const (
	vanilla iceCreamFlavor = iota
	chocolate
	strawberry
	peach
)

func (s iceCreamFlavor) String() string {
	return [...]string{"Vanilla", "Chocolate", "Strawberry", "Peach"}[s]
}

type person struct {
	firstName string
	lastName  string
	favFlavor iceCreamFlavor
}

func main() {
	people := []person{
		person{
			firstName: "James",
			lastName:  "Bond",
			favFlavor: vanilla,
		},
		person{
			firstName: "Miss",
			lastName:  "Moneypenny",
			favFlavor: chocolate,
		},
	}

	for _, v := range people {
		fmt.Printf("%v %v %v\n", v.firstName, v.lastName, v.favFlavor)
	}
}
