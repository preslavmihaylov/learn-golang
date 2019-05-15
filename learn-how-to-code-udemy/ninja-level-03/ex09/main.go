package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	favSport, _ := reader.ReadString('\n')
	favSport = strings.TrimRight(favSport, "\n")

	switch favSport {
	case "football":
		fmt.Println("Go Go Juventus")
	case "basketball":
		fmt.Println("Go Go Lakers")
	default:
		fmt.Println("I'm not a sports guy")
	}
}
