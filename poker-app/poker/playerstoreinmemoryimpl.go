package poker

import (
	"fmt"
)

type defaultStore struct {
	scoreboard map[string]int
}

func GetInMemoryStore(iniData ...map[string]int) PlayerStore {
	if len(iniData) > 0 {
		return &defaultStore{iniData[0]}
	}
	return &defaultStore{map[string]int{}}
}

func (s *defaultStore) GetScore(player string) (int, error) {
	score, ok := s.scoreboard[player]
	if ok {
		return score, nil
	}
	return 0, fmt.Errorf("Unknown player found %s", player)
}

func (s *defaultStore) RecordWin(name string) {
	s.scoreboard[name]++
}

func (s *defaultStore) GetLeagueTable() League {
	var result League
	for k, v := range s.scoreboard {
		result = append(result, Player{k, v})
	}
	return result
}
