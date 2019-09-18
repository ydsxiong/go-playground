package poker

import (
	"io"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, to io.Writer)
	Finish(winner string)
}

type pokerGame struct {
	store PlayerStore
	alert BlindAlerter
}

func NewGame(store PlayerStore, alert BlindAlerter) *pokerGame {
	return &pokerGame{store, alert}
}

func (g *pokerGame) Start(numberofplayers int, to io.Writer) {
	blindIncrement := time.Duration(5+numberofplayers) * time.Second
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alert.ScheduleAlertAt(blindTime, blind, to)
		blindTime = blindTime + blindIncrement
	}
}

func (g *pokerGame) Finish(winner string) {
	g.store.RecordWin(winner)
}
