package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	log.SetOutput(os.Stderr)

	os.Exit(code)
}

func TestResetCounting(t *testing.T) {
	c := BlackjackCounting{}
	c.runningCnt = 10
	c.discarded = 124
	c.totalDecks = 5
	c.minTrueCnt = 10

	c.reset()
	if c.runningCnt != 0 || c.discarded != 0 {
		t.Errorf("Expected runCount=%d, discarded=%d. Got runCount=%d, discarded=%d",
			0, 0, c.runningCnt, c.discarded)
	}

	if c.totalDecks != 5 || c.minTrueCnt != 10 {
		t.Errorf("Expected totalDecks=%d, minTrueCount=%d. Got totalDecks=%d, minTrueCount=%d",
			5, 10, c.totalDecks, c.minTrueCnt)
	}
}

func TestDeckShufflingResetsCounting(t *testing.T) {
	bs := NewStrategy(10, 5, 5, 5)
	bs.counting.runningCnt = 10
	bs.counting.discarded = 100
	bs.Listen(api.DeckShuffledEvent{})
	if bs.counting.runningCnt != 0 || bs.counting.discarded != 0 {
		t.Errorf("Expected counting to be reset after deck shuffling. Got runCount=%d, discarded=%d",
			bs.counting.runningCnt, bs.counting.discarded)
	}
}

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
			ExpectedRunningCnt: 2,
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
			ExpectedRunningCnt: 5,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.HitEvent{Card: cardWith(decks.Ace)},
				api.HitEvent{Card: cardWith(decks.King)},
				api.HitEvent{Card: cardWith(decks.Queen)},
				api.HitEvent{Card: cardWith(decks.Jack)},
				api.HitEvent{Card: cardWith(decks.Ten)},
			},
			ExpectedRunningCnt: -5,
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
			ExpectedRunningCnt: 5,
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
			ExpectedRunningCnt: -1,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.DealCardsEvent{
					DealerHand: handWith(decks.Two),
					Hands: map[string][]decks.Card{
						"Player 1": handWith(decks.Ten, decks.Four),
						"Player 2": handWith(decks.Jack, decks.Six),
					},
				},
				api.DealerTurnEvent{
					DealerRevealed: true,
					DealerHand:     handWith(decks.Two, decks.Three),
				},
			},
			ExpectedRunningCnt: 2,
		},
		RunCntData{
			Events: []api.GameEvent{
				api.DealCardsEvent{
					DealerHand: handWith(decks.Two),
					Hands: map[string][]decks.Card{
						"Player 1": handWith(decks.Ten, decks.Four),
						"Player 2": handWith(decks.Jack, decks.Six),
					},
				},
				api.DealerTurnEvent{
					DealerRevealed: false,
					DealerHand:     handWith(decks.Two, decks.Three),
				},
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
	tests := []TrueCntData{
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       0,
			GivenRunningCnt:      21,
			ExpectedBettingUnits: 7,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       52,
			GivenRunningCnt:      22,
			ExpectedBettingUnits: 11,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       104,
			GivenRunningCnt:      22,
			ExpectedBettingUnits: 22,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       27,
			GivenRunningCnt:      25,
			ExpectedBettingUnits: 10,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       52*2 + 27,
			GivenRunningCnt:      10,
			ExpectedBettingUnits: 20,
		},
		TrueCntData{
			GivenDecksCnt:        3,
			GivenDiscarded:       52*2 + 27,
			GivenRunningCnt:      -30,
			ExpectedBettingUnits: 0,
		},
		TrueCntData{
			GivenDecksCnt:        1,
			GivenDiscarded:       52,
			GivenRunningCnt:      30,
			ExpectedBettingUnits: 0,
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
