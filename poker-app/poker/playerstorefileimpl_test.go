package poker_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ydsxiong/go-playground/poker-app/poker"
)

func TestFileSystemStore(t *testing.T) {
	//strings.NewReader(str)
	initData := `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`

	database, err, clearDatabase := createTempFile(initData)
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer clearDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	if err != nil {
		t.Fatalf("Problem with opening file database, %v", err)
	}

	t.Run("/league from a reader", func(t *testing.T) {
		got := store.GetLeagueTable()
		want := poker.League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertleague(t, got, want)

		got = store.GetLeagueTable()
		assertleague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {

		got, _ := store.GetScore("Chris")
		want := 33

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		store.RecordWin("Chris")

		got, _ := store.GetScore("Chris")
		want := 34
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

		store.RecordWin("Cleo")

		got, _ = store.GetScore("Cleo")
		want = 11
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}

		store.RecordWin("Chris")

		got, _ = store.GetScore("Chris")
		want = 35
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("store wins for new player", func(t *testing.T) {
		store.RecordWin("Pepper")

		got, _ := store.GetScore("Pepper")
		want := 1
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, _, cleanDatabase := createTempFile("")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)

		if err != nil {
			t.Fatalf("Problem with opening file database, %v", err)
		}
	})
}

func TestTape_Write(t *testing.T) {
	file, err, clean := createTempFile("12345")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer clean()

	tape := poker.NewTapeFile(file)

	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func createTempFile(initialData string) (*os.File, error, func()) {

	tmpfile, err := ioutil.TempFile(".", "db")

	if err != nil {
		return nil, err, nil
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, nil, removeFile
}
