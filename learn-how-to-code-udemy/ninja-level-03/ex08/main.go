package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")

	switch {
	case text == "James":
		fmt.Println("This is James")
	case text == "Bond":
		fmt.Println("This is Bond")
	default:
		fmt.Println("I have no idea who you are")
	}
}
