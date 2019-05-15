package main

import "fmt"

func main() {
	arr1 := []string{"James", "Bond", "Shaken, not stirred"}
	arr2 := []string{"Miss", "Moneypenny", "Hellooooooo, James"}
	arr := [][]string{arr1, arr2}

	for _, arr1 := range arr {
		for _, v := range arr1 {
			fmt.Print(v, " ")
		}

		fmt.Println()
	}
}
