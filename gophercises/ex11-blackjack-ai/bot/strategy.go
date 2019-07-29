package bot

import (
	"log"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex09-deck/decks"
	"github.com/preslavmihaylov/learn-golang/gophercises/ex10-blackjack/blackjack/api"
)

func NewStrategy(roundsCnt int) *BlackjackStrategy {
	return &BlackjackStrategy{
		roundsCnt: roundsCnt,
	}
}

type BlackjackStrategy struct {
	Stats      BlackjackStats
	roundsCnt  int
	hasSplit   bool
	playerHand []decks.Card
	dealerHand []decks.Card
}

func (bb *BlackjackStrategy) Listen(e api.GameEvent) {
	switch ev := e.(type) {
	case api.StartBetEvent:
		bb.Stats.Balance = ev.Balance

		log.Printf("[StartBetEvent]: Balance=%d", ev.Balance)
	case api.PlayerTurnEvent:
		bb.playerHand = ev.PlayerHand
		bb.dealerHand = ev.DealerHand

		log.Printf("[PlayerTurnEvent]: Received %s %s", bb.playerHand, bb.dealerHand)
	case api.DealerTurnEvent:
		log.Printf("[DealerTurnEvent] Dealer's hand: %s", ev.DealerHand)
	case api.DoubleDownEvent:
		log.Printf("[DoubleDownEvent]: Received %s", ev.Card)
	case api.HitEvent:
		log.Printf("[HitEvent]: %s hits. Received %s", ev.PlayerName, ev.Card)
		if ev.Busted {
			log.Printf("[HitEvent]: %s busted!", ev.PlayerName)
		}
	case api.StandEvent:
		log.Printf("[StandEvent] %s stands", ev.PlayerName)
	case api.BlackjackEvent:
		log.Println("[PlayerBlackjack] Got Blackjack!")
	case api.ResolveEvent:
		for _, res := range ev.Results {
			if res.Outcome == api.PlayerBlackjack ||
				res.Outcome == api.DealerBusted ||
				res.Outcome == api.Won {

				bb.Stats.HandsWon++
			} else if res.Outcome == api.Tied {
				bb.Stats.HandsTied++
			} else {
				bb.Stats.HandsLost++
			}

			log.Printf("[ResolveEvent]: Outcome=%s", res.Outcome)
		}
	case api.RoundEndsEvent:
		bb.roundsCnt--
		bb.hasSplit = false

		log.Printf("[RoundEndsEvent]: %d rounds left", bb.roundsCnt)
	}
}

func (bb *BlackjackStrategy) BetTurn(actions []api.Action) api.Action {
	if bb.roundsCnt <= 0 {
		return &api.ExitAction{}
	}

	log.Println("[Bet Turn]: bet 100")
	return &api.BetAction{
		Bet: 100,
	}
}

func (bb *BlackjackStrategy) PlayerTurn(actions []api.Action) api.Action {
	playerScore := api.CalculateScore(bb.playerHand)

	var res api.Action
	if !bb.hasSplit && len(bb.playerHand) == 2 && bb.playerHand[0] == bb.playerHand[1] {
		bb.hasSplit = true
		res = bb.splitRules()
	} else if playerScore.IsSoft {
		res = bb.softRules()
	} else {
		res = bb.hardRules()
	}

	log.Printf("[PlayerTurn]: %s", res)
	return res
}

func (bb *BlackjackStrategy) splitRules() api.Action {
	pair := bb.playerHand[0].Rank
	dealerScore := api.CalculateScore(bb.dealerHand)
	switch pair {
	case decks.Ten:
		return &api.StandAction{}
	case decks.Nine:
		if bb.dealerHand[0].Rank != decks.Seven && dealerScore.Value >= 2 && dealerScore.Value < 10 {
			return &api.SplitAction{}
		}

		return &api.StandAction{}
	case decks.Eight:
		return &api.SplitAction{}
	case decks.Seven:
		if dealerScore.Value >= 2 && dealerScore.Value < 8 {
			return &api.SplitAction{}
		}

		return &api.HitAction{}
	case decks.Six:
		if dealerScore.Value >= 2 && dealerScore.Value < 7 {
			return &api.SplitAction{}
		}

		return &api.HitAction{}
	case decks.Five:
		if dealerScore.Value >= 2 && dealerScore.Value < 10 {
			return &api.DoubleAction{}
		}

		return &api.HitAction{}
	case decks.Four:
		if dealerScore.Value >= 5 && dealerScore.Value < 7 {
			return &api.SplitAction{}
		}

		return &api.HitAction{}
	case decks.Three:
		fallthrough
	case decks.Two:
		if dealerScore.Value >= 2 && dealerScore.Value < 8 {
			return &api.SplitAction{}
		}

		return &api.HitAction{}
	case decks.Ace:
		fallthrough
	default:
		return &api.SplitAction{}
	}
}

func (bb *BlackjackStrategy) softRules() api.Action {
	playerScore := api.CalculateScore(bb.playerHand)
	dealerScore := api.CalculateScore(bb.dealerHand)
	if playerScore.Value == 20 {
		return &api.StandAction{}
	} else if playerScore.Value == 19 {
		if len(bb.playerHand) == 2 && dealerScore.Value == 6 {
			return &api.DoubleAction{}
		}

		return &api.StandAction{}
	} else if playerScore.Value == 18 {
		if len(bb.playerHand) == 2 && dealerScore.Value >= 2 && dealerScore.Value < 7 {
			return &api.DoubleAction{}
		}

		if dealerScore.Value >= 9 && dealerScore.Value < 12 {
			return &api.HitAction{}
		}

		return &api.StandAction{}
	} else if playerScore.Value == 17 {
		if len(bb.playerHand) == 2 && dealerScore.Value >= 3 && dealerScore.Value < 7 {
			return &api.DoubleAction{}
		}
	} else if playerScore.Value >= 15 && playerScore.Value < 17 {
		if len(bb.playerHand) == 2 && dealerScore.Value >= 4 && dealerScore.Value < 7 {
			return &api.DoubleAction{}
		}
	} else if playerScore.Value >= 13 && playerScore.Value < 15 {
		if len(bb.playerHand) == 2 && dealerScore.Value >= 5 && dealerScore.Value < 7 {
			return &api.DoubleAction{}
		}
	}

	return &api.HitAction{}
}

func (bb *BlackjackStrategy) hardRules() api.Action {
	playerScore := api.CalculateScore(bb.playerHand)
	dealerScore := api.CalculateScore(bb.dealerHand)

	if playerScore.Value >= 17 {
		return &api.StandAction{}
	} else if playerScore.Value >= 13 && playerScore.Value < 17 {
		if dealerScore.Value >= 2 && dealerScore.Value < 7 {
			return &api.StandAction{}
		}
	} else if playerScore.Value == 12 {
		if dealerScore.Value >= 4 && dealerScore.Value < 7 {
			return &api.StandAction{}
		}
	}

	if len(bb.playerHand) != 2 {
		return &api.HitAction{}
	}

	if playerScore.Value == 11 {
		return &api.DoubleAction{}
	} else if playerScore.Value == 10 {
		if dealerScore.Value >= 2 && dealerScore.Value < 10 {
			return &api.DoubleAction{}
		}
	} else if playerScore.Value == 9 {
		if dealerScore.Value >= 3 && dealerScore.Value < 7 {
			return &api.DoubleAction{}
		}
	}

	return &api.HitAction{}
}
