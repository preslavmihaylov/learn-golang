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
	people := map[string]person{
		"Bond": person{
			firstName: "James",
			lastName:  "Bond",
			favFlavor: vanilla,
		},
		"Moneypenny": person{
			firstName: "Miss",
			lastName:  "Moneypenny",
			favFlavor: chocolate,
		},
	}

	for k, v := range people {
		fmt.Printf("%v -> %v %v %v\n", k, v.firstName, v.lastName, v.favFlavor)
	}
}
