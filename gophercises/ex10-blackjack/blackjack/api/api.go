package api

type BlackjackAPI interface {
	Listen(e GameEvent)
	BetTurn(actions []Action) Action
	PlayerTurn(actions []Action) Action
}
