package main

import (
	"fmt"
	"github.com/preslavmihaylov/learn-golang/ninja-level-13/ex02/quote"
	"github.com/preslavmihaylov/learn-golang/ninja-level-13/ex02/word"
)

func main() {
	fmt.Println(word.Count(quote.SunAlso))

	for k, v := range word.UseCount(quote.SunAlso) {
		fmt.Println(v, k)
	}
}
