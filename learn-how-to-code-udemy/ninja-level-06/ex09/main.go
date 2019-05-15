package main

import "fmt"

func main() {
	nums := []int{5, 2, 1, 10, 15, 12, 7}
	sort(func(a int, b int) int {
		return a - b
	}, &nums)

	fmt.Println(nums)
}

func sort(comp func(int, int) int, nums *[]int) {
	for i := range *nums {
		for j := 0; j < len(*nums)-i-1; j++ {
			if comp((*nums)[j], (*nums)[j+1]) > 0 {
				(*nums)[j], (*nums)[j+1] = (*nums)[j+1], (*nums)[j]
			}
		}
	}
}
