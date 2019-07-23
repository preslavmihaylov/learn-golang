package api

import (
	"fmt"
)

type Action interface {
	fmt.Stringer
	Help() string
}

func NewActions(actions ...Action) []Action {
	var res []Action
	for _, a := range actions {
		res = append(res, a)
	}

	res = append(res, ExitAction{})
	res = append(res, HelpAction{res})
	return res
}

type HelpAction struct {
	actions []Action
}

func (ha HelpAction) String() string {
	return "help"
}

func (ha HelpAction) Help() string {
	return "show info about available options"
}

type HitAction struct{}

func (ha HitAction) String() string {
	return "hit"
}

func (ha HitAction) Help() string {
	return "draw a new card"
}

type StandAction struct{}

func (ha StandAction) String() string {
	return "stand"
}

func (ha StandAction) Help() string {
	return "end turn and proceed with next player"
}

type ExitAction struct{}

func (e ExitAction) String() string {
	return "exit"
}

func (e ExitAction) Help() string {
	return "exit the game"
}
