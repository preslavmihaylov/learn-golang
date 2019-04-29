package main

import "fmt"

func main() {
	arr := []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51}

	arr = append(arr, 52)
	fmt.Println(arr)

	arr = append(arr, 53, 54, 55)
	fmt.Println(arr)

	arr = append(arr, []int{56, 57, 58, 59, 60}...)
	fmt.Println(arr)
}
