package blackjack

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/data"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/states"
)

func Play(players []data.PlayerInterface) {
	state := states.InitState
	gd := data.New(3, data.NewPlayer(players[0], "Player 1"))
	for state != nil {
		state = state(gd)
	}

	fmt.Println("Game Over")
}
