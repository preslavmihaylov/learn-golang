package main

import "fmt"

const (
	y1 = 2019 + iota
	y2 = 2019 + iota
	y3 = 2019 + iota
	y4 = 2019 + iota
)

func main() {
	fmt.Println(y1, y2, y3, y4)
}
