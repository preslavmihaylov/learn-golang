package blackjack

import (
	"fmt"

	bjapi "github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/internal/data"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/internal/states"
)

func Play(decksCnt, playersCnt int, api bjapi.BlackjackAPI) {
	state := states.InitState
	gd := data.New(decksCnt, playersCnt, api)
	for state != nil {
		state = state(gd)
	}

	fmt.Println("Game Over")
}
