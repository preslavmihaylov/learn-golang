package bot

import (
	"testing"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

func TestBalanceUpdateOnStartBet(t *testing.T) {
	bs := NewStrategy(5, 5, 5, 5)
	bs.Stats.Balance = 100
	bs.Listen(api.StartBetEvent{
		Balance: 50,
	})

	if bs.Stats.Balance != 50 {
		t.Errorf("Expected balance to update to 50. Got Balance=%d", bs.Stats.Balance)
	}
}

func TestHandsWonUpdated(t *testing.T) {
	bs := NewStrategy(5, 5, 5, 5)
	bs.Stats.HandsWon = 100
	expected := 100

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.Won,
			},
		},
	})

	expected++
	if bs.Stats.HandsWon != expected {
		t.Errorf("Expected hands won to update on win outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsWon)
	}

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.PlayerBlackjack,
			},
		},
	})

	expected++
	if bs.Stats.HandsWon != expected {
		t.Errorf("Expected hands won to update on player blackjack outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsWon)
	}

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.DealerBusted,
			},
		},
	})

	expected++
	if bs.Stats.HandsWon != expected {
		t.Errorf("Expected hands won to update on dealer busted outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsWon)
	}
}

func TestHandsTiedUpdated(t *testing.T) {
	bs := NewStrategy(5, 5, 5, 5)
	bs.Stats.HandsTied = 100
	expected := 100

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.Tied,
			},
		},
	})

	expected++
	if bs.Stats.HandsTied != expected {
		t.Errorf("Expected hands tied to update on tie outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsTied)
	}
}

func TestHandsLostUpdated(t *testing.T) {
	bs := NewStrategy(5, 5, 5, 5)
	bs.Stats.HandsLost = 100
	expected := 100

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.Lost,
			},
		},
	})

	expected++
	if bs.Stats.HandsLost != expected {
		t.Errorf("Expected hands lost to update on lost outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsLost)
	}

	bs.Listen(api.ResolveEvent{
		Results: map[string]api.Result{
			"Player 1": api.Result{
				Outcome: api.Busted,
			},
		},
	})

	expected++
	if bs.Stats.HandsLost != expected {
		t.Errorf("Expected hands lost to update on busted outcome."+
			" Expected %d, Got %d", expected, bs.Stats.HandsLost)
	}
}
