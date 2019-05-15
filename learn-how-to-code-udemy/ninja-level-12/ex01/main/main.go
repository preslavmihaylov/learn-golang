package main

import (
	"fmt"
	"github.com/preslavmihaylov/learn-golang/ninja-level-12/ex01/dog"
)

func main() {
	humanYears := 5
	dogYears := dog.Years(humanYears)
	fmt.Printf("%v human years == %v dog years\n", humanYears, dogYears)
}
