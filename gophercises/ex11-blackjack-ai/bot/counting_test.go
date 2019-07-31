package bot

import (
	"fmt"
	"testing"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

type RunCntData struct {
	Events             []api.GameEvent
	ExpectedRunningCnt int
}

func TestRunningCount(t *testing.T) {
	tests := []RunCntData{
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Two)},
				api.DoubleDownEvent{Card: cardWith(decks.Three)},
			},
			ExpectedRunningCnt: -2,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.DealCardsEvent{
					DealerHand: handWith(decks.Two),
					Hands: map[string][]decks.Card{
						"Player 1": handWith(decks.Three, decks.Four),
						"Player 2": handWith(decks.Five, decks.Six),
					},
				},
			},
			ExpectedRunningCnt: -5,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Ace)},
				api.HitEvent{Card: cardWith(decks.King)},
				api.HitEvent{Card: cardWith(decks.Queen)},
				api.HitEvent{Card: cardWith(decks.Jack)},
				api.HitEvent{Card: cardWith(decks.Ten)},
			},
			ExpectedRunningCnt: 5,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Seven)},
				api.HitEvent{Card: cardWith(decks.Eight)},
				api.HitEvent{Card: cardWith(decks.Nine)},
			},
			ExpectedRunningCnt: 0,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Two)},
				api.HitEvent{Card: cardWith(decks.Three)},
				api.HitEvent{Card: cardWith(decks.Four)},
				api.HitEvent{Card: cardWith(decks.Five)},
				api.HitEvent{Card: cardWith(decks.Six)},
			},
			ExpectedRunningCnt: -5,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Two)},
				api.HitEvent{Card: cardWith(decks.Three)},
				api.HitEvent{Card: cardWith(decks.Eight)},
				api.HitEvent{Card: cardWith(decks.Ten)},
				api.HitEvent{Card: cardWith(decks.Jack)},
				api.HitEvent{Card: cardWith(decks.Queen)},
			},
			ExpectedRunningCnt: 1,
		},
	}

	testRunning(t, tests)
}

func testRunning(t *testing.T, tests []RunCntData) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			str := NewStrategy(999, 3, 25, 4)
			for _, e := range test.Events {
				str.Listen(e)
			}

			if str.counting.runningCnt != test.ExpectedRunningCnt {
				t.Errorf("Expected %d, got %d", test.ExpectedRunningCnt, str.counting.runningCnt)
			}
		})
	}
}

type TrueCntData struct {
	GivenDecksCnt        int
	GivenDiscarded       int
	GivenRunningCnt      int
	ExpectedBettingUnits int
}

func TestTrueCnt(t *testing.T) {
	// TODO: write tests
	tests := []TrueCntData{
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       0,
			GivenRunningCnt:      21,
			ExpectedBettingUnits: 6,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       52,
			GivenRunningCnt:      22,
			ExpectedBettingUnits: 10,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       104,
			GivenRunningCnt:      22,
			ExpectedBettingUnits: 21,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       27,
			GivenRunningCnt:      25,
			ExpectedBettingUnits: 9,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       52*2 + 27,
			GivenRunningCnt:      10,
			ExpectedBettingUnits: 19,
		},
	}

	testTrueCnt(t, tests)
}

func testTrueCnt(t *testing.T, tests []TrueCntData) {
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			str := NewStrategy(999, test.GivenDecksCnt, 25, 4)
			str.counting.discarded = test.GivenDiscarded
			str.counting.runningCnt = test.GivenRunningCnt

			act := str.BetTurn([]api.Action{&api.BetAction{}}).(*api.BetAction)

			expectedBet := test.ExpectedBettingUnits * 25
			if expectedBet != act.Bet {
				t.Errorf("Expected %d, got %d", expectedBet, act.Bet)
			}
		})
	}
}

func cardWith(r decks.Rank) decks.Card {
	return decks.Card{
		Rank: r,
	}
}
