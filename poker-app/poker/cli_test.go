package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ydsxiong/go-playground/poker-app/poker"
)

var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type GameSpy struct {
	startedWith  int
	finishedWith string

	BlindAlert []byte
}

func (gs *GameSpy) Start(numberOfPlayers int, to io.Writer) {
	gs.startedWith = numberOfPlayers
	to.Write(gs.BlindAlert)
}
func (gs *GameSpy) Finish(winner string) {
	gs.finishedWith = winner
}

func TestCLI(t *testing.T) {

	db, err, clean := createTempFile("")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer clean()

	playerStore, e := poker.NewFileSystemPlayerStore(db)
	if e != nil {
		t.Errorf("problem initialising player db file, %v", e)
	}

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		cli := poker.NewPokerCLI(dummyStdIn, stdout, poker.NewGame(playerStore, dummySpyAlerter))
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt //"Please enter the number of players: "

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}

		out := &bytes.Buffer{}
		in := userSends("3", "Chris wins")

		poker.NewPokerCLI(in, out, game).PlayPoker()

		assertMessagesSentToUser(t, out, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayersWanted int) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.startedWith == numberOfPlayersWanted
	})
	if !passed {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.startedWith)
	}
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, winner string) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.finishedWith == winner
	})
	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.finishedWith)
	}
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func retryUntil(duration time.Duration, f func() bool) bool {
	deadline := time.Now().Add(duration)

	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
