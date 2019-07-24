package api

import (
	"fmt"
	"strconv"
)

type Action interface {
	fmt.Stringer
	Help() string
	ArgsCnt() int
	SetArgs(args ...string) error
}

func NewActions(actions ...Action) []Action {
	var res []Action
	for _, a := range actions {
		res = append(res, a)
	}

	res = append(res, &ExitAction{})
	res = append(res, &HelpAction{NoArgsAction{}, res})
	return res
}

type NoArgsAction struct{}

func (_ NoArgsAction) ArgsCnt() int {
	return 0
}

func (_ NoArgsAction) SetArgs(args ...string) error {
	if len(args) > 0 {
		return fmt.Errorf("Too many args passed")
	}

	return nil
}

type HelpAction struct {
	NoArgsAction
	actions []Action
}

func (ha *HelpAction) String() string {
	return "help"
}

func (ha *HelpAction) Help() string {
	return "show info about available options"
}

type BetAction struct {
	Bet int
}

func (ba *BetAction) String() string {
	return "bet"
}

func (ba *BetAction) Help() string {
	return "place a bet"
}

func (ba *BetAction) ArgsCnt() int {
	return 1
}

func (ba *BetAction) SetArgs(args ...string) error {
	if len(args) > 1 {
		return fmt.Errorf("Too many args passed")
	}

	var err error
	ba.Bet, err = strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse bet: %s", err)
	}

	return nil
}

type HitAction struct {
	NoArgsAction
}

func (ha *HitAction) String() string {
	return "hit"
}

func (ha *HitAction) Help() string {
	return "draw a new card"
}

type StandAction struct {
	NoArgsAction
}

func (ha *StandAction) String() string {
	return "stand"
}

func (ha *StandAction) Help() string {
	return "end turn and proceed with next player"
}

type ExitAction struct {
	NoArgsAction
}

func (e *ExitAction) String() string {
	return "exit"
}

func (e *ExitAction) Help() string {
	return "exit the game"
}
