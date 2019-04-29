package main

import "fmt"

func main() {
	arr := []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51}

	arr1 := arr[0:5]
	arr2 := arr[5:]
	arr3 := arr[2:7]
	arr4 := arr[1:6]

	fmt.Println(arr1)
	fmt.Println(arr2)
	fmt.Println(arr3)
	fmt.Println(arr4)
}
