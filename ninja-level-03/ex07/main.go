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

	if text == "James" {
		fmt.Println("This is James")
	} else if text == "Bond" {
		fmt.Println("This is Bond")
	} else {
		fmt.Println("I have no idea who you are")
	}
}
