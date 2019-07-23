package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type CLIPlayer struct{}

func (c *CLIPlayer) Listen(ev api.GameEvent) {
	switch e := ev.(type) {
	case api.DealCardsEvent:
		fmt.Println("--- Dealing Cards ---")
		for name, hand := range e.Hands {
			fmt.Printf("%s's hand:\n", name)
			for _, c := range hand {
				fmt.Printf("\t%s\n", c)
			}
		}

		fmt.Printf("Dealer's hand:")
		for _, c := range e.DealerHand {
			fmt.Printf("\t%s\n", c)
		}
	case api.PlayerTurnEvent:
		fmt.Printf("--- %s's turn ---\n", e.PlayerName)
		fmt.Println("Player's Hand:")
		for _, c := range e.PlayerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Println()
		fmt.Printf("%s's Score: %d\n", e.PlayerName, api.CalculateScore(e.PlayerHand).Value)

		fmt.Println("Dealer's Hand:")
		for _, c := range e.DealerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Println()
		fmt.Printf("Dealer's Score: %d\n", api.CalculateScore(e.DealerHand).Value)
	case api.DealerTurnEvent:
		fmt.Println("--- Dealer's turn ---")
		fmt.Println("Players still in game:")
		for name, hand := range e.PlayersInGame {
			fmt.Printf("%s's Hand:\n", name)
			for _, c := range hand {
				fmt.Printf("\t%s\n", c)
			}

			fmt.Printf("\n%s's Score: %d\n", name, api.CalculateScore(hand).Value)
		}

		if e.DealerRevealed {
			fmt.Println("Dealer reveals hand...")
		}

		fmt.Println("Dealer's Hand:")
		for _, c := range e.DealerHand {
			fmt.Printf("\t%s\n", c)
		}

		fmt.Printf("\nDealer's Score: %d\n", api.CalculateScore(e.DealerHand).Value)
	case api.HitEvent:
		fmt.Printf("--- %s hits! ---\n", e.PlayerName)
		fmt.Printf("Got %s\n", e.Card)
		if e.Busted {
			fmt.Printf("%s Busted!\n", e.PlayerName)
		}
	case api.StandEvent:
		fmt.Printf("--- %s stands. ---\n", e.PlayerName)
	case api.ResolveEvent:
		fmt.Println("--- Resolution ---")

		for name, res := range e.Results {
			if res.Outcome == api.Busted {
				fmt.Printf("%s Busted!\n", name)
				continue
			} else if res.Outcome == api.DealerBusted {
				fmt.Printf("Dealer Busted! %s won!\n", name)
				continue
			}

			fmt.Printf("%s's Score: %d, Dealer's Score: %d\n", name, res.PlayerScore, res.DealerScore)
			if res.Outcome == api.Won {
				fmt.Printf("%s won!\n", name)
			} else if res.Outcome == api.Tied {
				fmt.Printf("%s tied!\n", name)
			} else if res.Outcome == api.Lost {
				fmt.Printf("%s lost!\n", name)
			}
		}
	case api.RoundEndsEvent:
		fmt.Println("--- Round Ends ---")
	}

}

func (c *CLIPlayer) PlayerTurn(actions []api.Action) api.Action {
	fmt.Println("What will you do?")
	fmt.Printf("> ")

	var choice string
	fmt.Scanln(&choice)
	choice = strings.Trim(choice, " ")

	for _, a := range actions {
		if a.String() == choice {
			return a
		}
	}

	return nil
}

func main() {
	playersCnt := flag.Int("players", 1, "number of players in game")
	flag.Parse()

	blackjack.Play(*playersCnt, &CLIPlayer{})
}
