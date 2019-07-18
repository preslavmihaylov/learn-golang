package main

import (
	"fmt"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/data"
)

type CLIPlayer struct{}

func (c *CLIPlayer) PlayerTurn(gi data.GameInfo, actions []data.Action) data.Action {
	fmt.Println("What will you do?")
	fmt.Printf("> ")

	var choice string
	fmt.Scanf("%s", &choice)
	choice = strings.Trim(choice, " ")

	for _, a := range actions {
		if a.String() == choice {
			return a
		}
	}

	return nil
}

func main() {
	// TODO: Implement more than one player in game
	blackjack.Play([]data.PlayerInterface{&CLIPlayer{}})
}
