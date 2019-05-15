package main

import "fmt"

func main() {
	mp := map[string][]string{
		"Bond James":      {"Shaken, not stirred", "Martinis", "Women"},
		"Moneypenny Miss": {"James Bond", "Literature", "Computer Science"},
		"No Doctor":       {"Being Evil", "Ice Cream", "Sunsets"},
	}

	mp["Mcleod Todd"] = []string{"Go", "Talking", "Teaching"}
	delete(mp, "Bond James")

	for k, v := range mp {
		fmt.Println(k)
		for i, v1 := range v {
			fmt.Println("\t", i, v1)
		}
	}
}
