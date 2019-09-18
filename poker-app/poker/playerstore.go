package poker

type Player struct {
	Name string
	Wins int
}

type League []Player

func (l League) Find(name string) *Player {
	
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

type PlayerStore interface {
	GetScore(name string) (int, error)
	RecordWin(name string)
	GetLeagueTable() League
}
