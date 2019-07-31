package states

import (
	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/internal/data"
)

type GameState func(data *data.GameData) GameState

func InitState(gd *data.GameData) GameState {
	return transition(betState)
}

func betState(gd *data.GameData) GameState {
	if gd.IsDealersTurn() {
		gd.NextPlayersTurn()
		return transition(dealState)
	}

	pl := gd.CurrentPlayer()
	gd.API().Listen(api.StartBetEvent{
		PlayerName: pl.Name(),
		Balance:    pl.Balance(),
	})

	return betActions(gd)
}

func dealState(gd *data.GameData) GameState {
	var ev api.DealCardsEvent
	ev.Hands = make(map[string][]decks.Card)

	for i := 0; i < data.HandSize; i++ {
		for j, p := range gd.Players() {
			c := gd.Draw()
			gd.Players()[j].Deal(c)
			ev.Hands[p.Name()] = append(ev.Hands[p.Name()], c)
		}

		c := gd.Draw()
		gd.Dealer.Deal(c)
		ev.DealerHand = append(ev.DealerHand, c)
	}

	gd.API().Listen(ev)

	return transition(playerTurnState)
}

func playerTurnState(gd *data.GameData) GameState {
	if gd.IsDealersTurn() {
		return transition(dealerTurnState)
	}

	pl := gd.CurrentPlayer()
	gd.API().Listen(api.PlayerTurnEvent{
		PlayerName: pl.Name(),
		PlayerHand: pl.Hand(),
		DealerHand: gd.Dealer.Hand(),
	})

	if pl.HasBlackjack() {
		gd.API().Listen(api.BlackjackEvent{
			PlayerName: pl.Name(),
		})

		gd.NextPlayersTurn()
		return transition(playerTurnState)
	}

	return playerActions(gd)
}

func dealerTurnState(gd *data.GameData) GameState {
	var ev api.DealerTurnEvent
	ev.PlayersInGame = make(map[string][]decks.Card)

	for _, p := range gd.Players() {
		if !p.IsBusted() {
			ev.PlayersInGame[p.Name()] = p.Hand()
		}
	}

	if !gd.Dealer.Revealed() {
		ev.DealerRevealed = true
		gd.Dealer.Reveal()
	}

	ev.DealerHand = gd.Dealer.Hand()

	gd.API().Listen(ev)
	if dealerShouldPlay(gd) {
		return transition(hitState)
	}

	return transition(standState)
}

func resolveState(gd *data.GameData) GameState {
	var ev api.ResolveEvent
	ev.Results = make(map[string]api.Result)

	for _, p := range gd.Players() {
		res := api.Result{
			PlayerScore: p.Score(),
			DealerScore: gd.Dealer.Score(),
		}

		if p.IsBusted() {
			res.Outcome = api.Busted
		} else if gd.Dealer.IsBusted() {
			res.Outcome = api.DealerBusted
		} else if p.HasBlackjack() {
			res.Outcome = api.PlayerBlackjack
		} else if gd.Dealer.Score() < p.Score() {
			res.Outcome = api.Won
		} else if gd.Dealer.Score() > p.Score() {
			res.Outcome = api.Lost
		} else {
			res.Outcome = api.Tied
		}

		if res.Outcome == api.Won || res.Outcome == api.DealerBusted {
			p.Payout(1)
		} else if res.Outcome == api.PlayerBlackjack && p.IsSplit() {
			p.Payout(1)
		} else if res.Outcome == api.PlayerBlackjack {
			p.Payout(1.5)
		} else if res.Outcome == api.Tied {
			p.Payout(0)
		} else {
			p.LoseBet()
		}

		ev.Results[p.Name()] = res
	}

	gd.API().Listen(ev)
	return transition(roundEndsState)
}

func roundEndsState(gd *data.GameData) GameState {
	gd.API().Listen(api.RoundEndsEvent{})
	gd.NewRound()
	if gd.ShouldShuffle() {
		gd.Shuffle()
		gd.API().Listen(api.DeckShuffledEvent{})
	}

	return transition(betState)
}

func hitState(gd *data.GameData) GameState {
	var ev api.HitEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	c := gd.Draw()
	player.Deal(c)
	ev.Card = c

	if player.IsBusted() {
		ev.Busted = true
	}

	gd.API().Listen(ev)

	var nextState GameState
	if !gd.IsDealersTurn() {
		nextState = playerTurnState
		if player.IsBusted() {
			gd.NextPlayersTurn()
		}
	} else {
		nextState = dealerTurnState
		if player.IsBusted() {
			nextState = resolveState
		}
	}

	return transition(nextState)
}

func standState(gd *data.GameData) GameState {
	var ev api.StandEvent

	player := gd.CurrentPlayer()
	ev.PlayerName = player.Name()

	gd.API().Listen(ev)
	var nextState GameState
	if !gd.IsDealersTurn() {
		nextState = playerTurnState
	} else {
		nextState = resolveState
	}

	if !gd.IsDealersTurn() {
		gd.NextPlayersTurn()
	}

	return transition(nextState)
}

func exitState(gd *data.GameData) GameState {
	return nil
}
