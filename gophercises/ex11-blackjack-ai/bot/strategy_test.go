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

func TestStoppingSimulation(t *testing.T) {
	bs := NewStrategy(0, 5, 5, 5)
	act := bs.BetTurn([]api.Action{&api.BetAction{}, &api.HelpAction{}, &api.ExitAction{}})
	switch act.(type) {
	case *api.ExitAction:
		// do nothing
	default:
		t.Errorf("Expected exit action, got %s", act)
	}
}

func TestRoundDecrementing(t *testing.T) {
	bs := NewStrategy(10, 5, 5, 5)
	bs.Listen(api.RoundEndsEvent{})
	if bs.roundsCnt != 9 {
		t.Errorf("Expected rounds to decrement from 10 to 9 upon round end. Got %d", bs.roundsCnt)
	}
}

func TestHardTotals(t *testing.T) {
	tests := []TurnData{
		// score >= 17
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

		// score <= 11 && len(hand) > 2
		TurnData{handWith(decks.Seven, decks.Two, decks.Two), handWith(decks.Five), &api.HitAction{}},

		// Invalid length for double
		TurnData{handWith(decks.Jack), handWith(decks.Four), &api.HitAction{}},
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

func TestSplitRules(t *testing.T) {
	tests := []TurnData{
		// Tens
		TurnData{handWith(decks.Ten, decks.Ten), handWith(decks.Five), &api.StandAction{}},

		// Nines
		TurnData{handWith(decks.Nine, decks.Nine), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Nine, decks.Nine), handWith(decks.Nine), &api.SplitAction{}},
		TurnData{handWith(decks.Nine, decks.Nine), handWith(decks.Ten), &api.StandAction{}},

		// Eights
		TurnData{handWith(decks.Eight, decks.Eight), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Eight, decks.Eight), handWith(decks.Ten), &api.SplitAction{}},

		// Sevens
		TurnData{handWith(decks.Seven, decks.Seven), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Seven, decks.Seven), handWith(decks.Seven), &api.SplitAction{}},
		TurnData{handWith(decks.Seven, decks.Seven), handWith(decks.Eight), &api.HitAction{}},

		// Sixs
		TurnData{handWith(decks.Six, decks.Six), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Six, decks.Six), handWith(decks.Six), &api.SplitAction{}},
		TurnData{handWith(decks.Six, decks.Six), handWith(decks.Seven), &api.HitAction{}},

		// Fives
		TurnData{handWith(decks.Five, decks.Five), handWith(decks.Two), &api.DoubleAction{}},
		TurnData{handWith(decks.Five, decks.Five), handWith(decks.Nine), &api.DoubleAction{}},
		TurnData{handWith(decks.Five, decks.Five), handWith(decks.Ten), &api.HitAction{}},

		// Fours
		TurnData{handWith(decks.Four, decks.Four), handWith(decks.Five), &api.SplitAction{}},
		TurnData{handWith(decks.Four, decks.Four), handWith(decks.Six), &api.SplitAction{}},
		TurnData{handWith(decks.Four, decks.Four), handWith(decks.Four), &api.HitAction{}},

		// Threes
		TurnData{handWith(decks.Three, decks.Three), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Three, decks.Three), handWith(decks.Seven), &api.SplitAction{}},
		TurnData{handWith(decks.Three, decks.Three), handWith(decks.Eight), &api.HitAction{}},

		// Twos
		TurnData{handWith(decks.Two, decks.Two), handWith(decks.Two), &api.SplitAction{}},
		TurnData{handWith(decks.Two, decks.Two), handWith(decks.Seven), &api.SplitAction{}},
		TurnData{handWith(decks.Two, decks.Two), handWith(decks.Eight), &api.HitAction{}},

		// Aces
		TurnData{handWith(decks.Ace, decks.Ace), handWith(decks.Ace), &api.SplitAction{}},
		TurnData{handWith(decks.Ace, decks.Ace), handWith(decks.Ten), &api.SplitAction{}},

		// Invalid hand for split
		TurnData{handWith(decks.Ace, decks.Ace, decks.Two), handWith(decks.Ten), &api.HitAction{}},
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
