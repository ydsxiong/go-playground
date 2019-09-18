package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

type Tape struct {
	file *os.File
}

func (t *Tape) Write(b []byte) (int, error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(b)
}

func NewTapeFile(database *os.File) *Tape {
	return &Tape{database}
}

type FileSystemPlayerStore struct {
	league   League
	database *json.Encoder
}

func LoadUpFileStore(filestorepath string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(filestorepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", filestorepath, err)
	}
	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)
	if err != nil {
		return nil, nil, fmt.Errorf("Problem with opening file database, %v", err)
	}

	return store, closeFunc, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	e := initialisePlayerDBFile(file)
	if e != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", e)
	}

	league, err := newLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	store := &FileSystemPlayerStore{
		league:   league,
		database: json.NewEncoder(&Tape{file})}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for _ = range ticker.C {
			league, err := newLeague(file)
			if err == nil {
				store.league = league
			}
		}
	}()
	return store, nil
}

func (fs *FileSystemPlayerStore) GetLeagueTable() League {
	return fs.league
}

func (fs *FileSystemPlayerStore) GetScore(name string) (int, error) {
	player := fs.league.Find(name)

	if player != nil {
		return player.Wins, nil
	}

	return 0, fmt.Errorf("Unknown player %v", name)
}

/**
When you range over a slice you are returned the current index of the loop (in our case i)
and a copy of the element at that index. Changing the Wins value of a copy won't have any effect
on the league slice that we iterate on. For that reason, we need to get the reference to the actual value
 by doing league[i] and then changing that value instead.
*/
func (fs *FileSystemPlayerStore) RecordWin(name string) {
	player := fs.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		fs.league = append(fs.league, Player{name, 1})
	}
	sortLeague(fs.league)
	fs.database.Encode(fs.league)
}

func sortLeague(league League) {
	sort.Slice(league, func(i, j int) bool {
		samescore := league[i].Wins == league[j].Wins
		if samescore {
			return league[i].Name < league[j].Name
		}
		return league[i].Wins > league[j].Wins
	})
}

func newLeague(data io.ReadSeeker) (l League, err error) {
	data.Seek(0, 0)
	var p []Player
	err = json.NewDecoder(data).Decode(&p)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
		l = nil
	} else {
		l = League(p)
		sortLeague(l)
	}
	return
}
