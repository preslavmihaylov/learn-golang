package main

import "fmt"

func main() {
	inc1 := incrementor()
	inc2 := incrementor()

	fmt.Println(inc1())
	fmt.Println(inc1())
	fmt.Println(inc1())

	fmt.Println(inc2())
	fmt.Println(inc2())
	fmt.Println(inc2())
}

func incrementor() func() int {
	x := 0
	return func() int {
		x++
		return x
	}
}
