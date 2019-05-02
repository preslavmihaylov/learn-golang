package main

import (
	"fmt"
	"sort"
)

type user struct {
	First   string
	Last    string
	Age     int
	Sayings []string
}

// ByAgeThenLast user sort comparator by age, then by last name
type ByAgeThenLast []user

func (u ByAgeThenLast) Len() int      { return len(u) }
func (u ByAgeThenLast) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
func (u ByAgeThenLast) Less(i, j int) bool {
	if u[i].Age == u[j].Age {
		return u[i].Last < u[j].Last
	}

	return u[i].Age < u[j].Age
}

func main() {
	u1 := user{
		First: "James",
		Last:  "Bond",
		Age:   32,
		Sayings: []string{
			"Shaken, not stirred",
			"Youth is no guarantee of innovation",
			"In his majesty's royal service",
		},
	}

	u2 := user{
		First: "Miss",
		Last:  "Moneypenny",
		Age:   32,
		Sayings: []string{
			"James, it is soo good to see you",
			"Would you like me to take care of that for you, James?",
			"I would really prefer to be a secret agent myself.",
		},
	}

	u3 := user{
		First: "M",
		Last:  "Hmmmm",
		Age:   54,
		Sayings: []string{
			"Oh, James. You didn't.",
			"Dear God, what has James done now?",
			"Can someone please tell me where James Bond is?",
		},
	}

	users := []user{u1, u2, u3}

	fmt.Println(users)
	fmt.Println()
	sort.Sort(ByAgeThenLast(users))
	for i, v := range users {
		sort.Strings(v.Sayings)
		fmt.Printf("User #%d\n", i)
		fmt.Println("\t", v.First, v.Last, v.Age)
		for j, s := range v.Sayings {
			fmt.Println("\t\t", j, s)
		}
	}

	// your code goes here

}
