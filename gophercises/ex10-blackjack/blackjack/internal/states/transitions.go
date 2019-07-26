package states

import "time"

func delayedTransition(gs GameState) GameState {
	time.Sleep(time.Second * 2)
	return gs
}

func transition(gs GameState) GameState {
	return gs
}
