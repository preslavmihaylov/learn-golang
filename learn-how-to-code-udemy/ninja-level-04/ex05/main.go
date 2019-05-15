package main

import "fmt"

func main() {
	arr := []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51}

	arr = append(arr[:3], arr[6:]...)
	fmt.Println(arr)
}
