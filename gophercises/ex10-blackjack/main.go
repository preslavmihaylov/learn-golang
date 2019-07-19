package main

import (
	"fmt"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/data"
)

type CLIPlayer struct{}

func (c *CLIPlayer) Listen(ev data.GameEvent) {
	switch e := ev.(type) {
	case data.DealCardsEvent:
		fmt.Println("--- Dealing Cards ---")
		_ = e
	case data.PlayerTurnEvent:
		fmt.Printf("--- %s's turn ---\n", e.PlayerName)
		fmt.Println("Player's Hand:")
		for _, c := range e.PlayerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Println()
		fmt.Printf("%s's Score: %d\n", e.PlayerName, data.CalculateScore(e.PlayerHand).Value)

		fmt.Println("Dealer's Hand:")
		for _, c := range e.DealerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Println()
		fmt.Printf("Dealer's Score: %d\n", data.CalculateScore(e.DealerHand).Value)
	case data.DealerTurnEvent:
		fmt.Println("--- Dealer's turn ---")
		fmt.Println("Players still in game:")
		for name, hand := range e.PlayersInGame {
			fmt.Printf("%s's Hand:\n", name)
			for _, c := range hand {
				fmt.Printf("\t%s\n", c)
			}

			fmt.Printf("\n%s's Score: %d\n", name, data.CalculateScore(hand).Value)
		}

		if e.DealerRevealed {
			fmt.Println("Dealer reveals hand...")
		}

		fmt.Println("Dealer's Hand:")
		for _, c := range e.DealerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Printf("\nDealer's Score: %d\n", data.CalculateScore(e.DealerHand).Value)
	case data.HitEvent:
		fmt.Printf("--- %s hits! ---\n", e.PlayerName)
		fmt.Printf("Got %s\n", e.Card)
		if e.Busted {
			fmt.Printf("%s Busted!\n", e.PlayerName)
		}
	case data.StandEvent:
		fmt.Printf("--- %s stands. ---\n", e.PlayerName)
	case data.ResolveEvent:
		fmt.Println("--- Resolution ---")

		name := "Player 1"
		res := e.Results["Player 1"]
		if res.Outcome == data.Busted {
			fmt.Printf("%s Busted!\n", name)
			break
		} else if res.Outcome == data.DealerBusted {
			fmt.Printf("Dealer Busted!\n")
			break
		}

		fmt.Printf("%s's Score: %d, Dealer's Score: %d\n", name, res.PlayerScore, res.DealerScore)
		if res.Outcome == data.Won {
			fmt.Printf("%s won!\n", name)
		} else if res.Outcome == data.Tied {
			fmt.Printf("%s tied!\n", name)
		} else if res.Outcome == data.Lost {
			fmt.Printf("%s lost!\n", name)
		}
	case data.RoundEndsEvent:
		fmt.Println("--- Round Ends ---")
	}

}

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
	// TODO: Rename PlayerInterface to Blackjack API
	// TODO: separate Blackjack API from blackjack/data package
	// TODO: separate Blackjack actions from blackjack/data package
	// TODO: make blackjack/data package part of blackjack/internal dir tree
	// TODO: make blackjack/states package part of blackjack/internal dir tree
	// TODO: Implement more than one player in game
	// TODO: DealCardsEvent handle more than one player
	blackjack.Play([]data.PlayerInterface{&CLIPlayer{}})
}
