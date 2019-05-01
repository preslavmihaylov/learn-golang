package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(foo(nums...))
	fmt.Println(bar(nums))
}

func foo(nums ...int) int {
	total := 0
	for _, v := range nums {
		total += v
	}

	return total
}

func bar(nums []int) int {
	return foo(nums...)
}
