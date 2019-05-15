package main

import "fmt"

func main() {
	x := 42
	y := 50
	equals := x == y
	le := x <= y
	ge := x >= y
	not := x != y
	lt := x < y
	gt := x > y

	fmt.Println("x == y", equals)
	fmt.Println("x <= y", le)
	fmt.Println("x >= y", ge)
	fmt.Println("x != y", not)
	fmt.Println("x < y", lt)
	fmt.Println("x > y", gt)
}
