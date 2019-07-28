package bot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type TurnData struct {
	PlayerHand     []decks.Card
	DealerHand     []decks.Card
	ExpectedAction api.Action
}

func TestHardTotals(t *testing.T) {
	tests := []TurnData{
		// score >= 17
		TurnData{handWith(decks.Ten, decks.Ten), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Nine), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Eight), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Seven), handWith(decks.Five), &api.StandAction{}},

		// score >= 13 && score < 17
		TurnData{handWith(decks.Ten, decks.Six), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Six), handWith(decks.Seven), &api.HitAction{}},
		TurnData{handWith(decks.Ten, decks.Five), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Five), handWith(decks.Seven), &api.HitAction{}},
		TurnData{handWith(decks.Ten, decks.Four), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Four), handWith(decks.Seven), &api.HitAction{}},
		TurnData{handWith(decks.Ten, decks.Three), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Three), handWith(decks.Seven), &api.HitAction{}},

		// score == 12
		TurnData{handWith(decks.Ten, decks.Two), handWith(decks.Five), &api.StandAction{}},
		TurnData{handWith(decks.Ten, decks.Two), handWith(decks.Three), &api.HitAction{}},

		// score == 11
		TurnData{handWith(decks.Nine, decks.Two), handWith(decks.Five), &api.DoubleAction{}},

		// score == 10
		TurnData{handWith(decks.Eight, decks.Two), handWith(decks.Seven), &api.DoubleAction{}},
		TurnData{handWith(decks.Eight, decks.Two), handWith(decks.Ten), &api.HitAction{}},

		// score == 9
		TurnData{handWith(decks.Seven, decks.Two), handWith(decks.Five), &api.DoubleAction{}},
		TurnData{handWith(decks.Seven, decks.Two), handWith(decks.Seven), &api.HitAction{}},

		// score <= 8
		TurnData{handWith(decks.Six, decks.Two), handWith(decks.Five), &api.HitAction{}},
		TurnData{handWith(decks.Five, decks.Two), handWith(decks.Five), &api.HitAction{}},
		TurnData{handWith(decks.Four, decks.Two), handWith(decks.Five), &api.HitAction{}},
		TurnData{handWith(decks.Three, decks.Two), handWith(decks.Five), &api.HitAction{}},
		TurnData{handWith(decks.Two, decks.Two), handWith(decks.Five), &api.HitAction{}},

		// score <= 11 && len(hand) > 2
		TurnData{handWith(decks.Seven, decks.Two, decks.Two), handWith(decks.Five), &api.HitAction{}},
	}

	turnTest(t, tests)
}

func TestSoftTotals(t *testing.T) {
	tests := []TurnData{
		// score == 20
		TurnData{handWith(decks.Ace, decks.Nine), handWith(decks.Five), &api.StandAction{}},

		// score == 19
		TurnData{handWith(decks.Ace, decks.Eight), handWith(decks.Six), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Eight), handWith(decks.Seven), &api.StandAction{}},
		TurnData{handWith(decks.Ace, decks.Six, decks.Two), handWith(decks.Six), &api.StandAction{}},

		// score == 18
		TurnData{handWith(decks.Ace, decks.Seven), handWith(decks.Six), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Seven), handWith(decks.Ace), &api.HitAction{}},
		TurnData{handWith(decks.Ace, decks.Seven), handWith(decks.Eight), &api.StandAction{}},
		TurnData{handWith(decks.Ace, decks.Five, decks.Two), handWith(decks.Six), &api.StandAction{}},

		// score == 17
		TurnData{handWith(decks.Ace, decks.Six), handWith(decks.Six), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Six), handWith(decks.Seven), &api.HitAction{}},

		// score >= 15 && score < 17
		TurnData{handWith(decks.Ace, decks.Five), handWith(decks.Six), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Four), handWith(decks.Six), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Five), handWith(decks.Seven), &api.HitAction{}},
		TurnData{handWith(decks.Ace, decks.Four), handWith(decks.Seven), &api.HitAction{}},
		TurnData{handWith(decks.Ace, decks.Two, decks.Two), handWith(decks.Six), &api.HitAction{}},

		// score >= 13 && score < 15
		TurnData{handWith(decks.Ace, decks.Three), handWith(decks.Five), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Two), handWith(decks.Five), &api.DoubleAction{}},
		TurnData{handWith(decks.Ace, decks.Three), handWith(decks.Four), &api.HitAction{}},
		TurnData{handWith(decks.Ace, decks.Two), handWith(decks.Four), &api.HitAction{}},
		TurnData{handWith(decks.Ace, decks.Ace, decks.Two), handWith(decks.Five), &api.HitAction{}},
	}

	turnTest(t, tests)
}

func turnTest(t *testing.T, tests []TurnData) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			ai := BlackjackStrategy{}
			ai.Listen(api.PlayerTurnEvent{
				PlayerName: "Player 1",
				PlayerHand: test.PlayerHand,
				DealerHand: test.DealerHand,
			})

			res := ai.PlayerTurn(nil)
			if reflect.TypeOf(res) != reflect.TypeOf(test.ExpectedAction) {
				t.Errorf("\nGiven Player=%s, Dealer=%s.\nExpected %s, got %s",
					test.PlayerHand, test.DealerHand, test.ExpectedAction, res)
			}
		})
	}
}

func handWith(ranks ...decks.Rank) []decks.Card {
	cards := []decks.Card{}
	for _, r := range ranks {
		cards = append(cards, decks.Card{Rank: r, Suit: decks.Clovers})
	}

	return cards
}
