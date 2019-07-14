package main

import (
	"fmt"
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
)

type Player struct {
	hand []decks.Card
}

type Dealer Player

type BlackjackData struct {
	deck    *decks.Deck
	players []Player
	dealer  Dealer
}

type BlackjackState func(data *BlackjackData) BlackjackState

func InitState(data *BlackjackData) BlackjackState {
	var err error
	data.deck, err = decks.New(decks.WithDecks(3), decks.Shuffle())
	if err != nil {
		log.Fatalf("failed to initialize deck: %s", err)
	}

	data.players = []Player{Player{}}
	data.dealer = Dealer{}

	return DealState
}

func DealState(data *BlackjackData) BlackjackState {
	fmt.Println("Dealing cards...")
	handSize := 2
	for i := 0; i < handSize; i++ {
		for i := range data.players {
			data.players[i].hand = append(data.players[i].hand, data.deck.Draw())
		}

		data.dealer.hand = append(data.dealer.hand, data.deck.Draw())
	}

	return PlayerTurn(0)
}

func PlayerTurn(i int) BlackjackState {
	return func(data *BlackjackData) BlackjackState {
		player := data.players[i]
		fmt.Printf("Player %d's turn\n", i)
		fmt.Println("Your hand:")
		for _, c := range player.hand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Println()
		fmt.Println("Dealer's hand:")
		fmt.Printf("\t%s\n", data.dealer.hand[0])

		fmt.Println()
		// your score
		// dealer score

		return nil
	}
}

func main() {
	state := InitState
	data := BlackjackData{}
	for state != nil {
		state = state(&data)
	}

	fmt.Println("Game Over")
}
