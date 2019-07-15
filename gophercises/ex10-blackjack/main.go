package main

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"
)

func main() {
	// TODO: More intelligent hit/stand decision for dealer:
	//		- is dealer score bigger than all players' score?
	//		- are all players busted?
	// TODO: ResolveState should check if dealer is busted
	// TODO: IsSoftScore is not working correctly.
	//		- score should include a boolean IsSoft in struct
	state := blackjack.InitState
	data := blackjack.BlackjackData{}
	for state != nil {
		state = state(&data)
	}

	fmt.Println("Game Over")
}
