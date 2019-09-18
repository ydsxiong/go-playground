package poker_test

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/ydsxiong/go-playground/poker-app/poker"
)

type scheduleAlert struct {
	scheduledAt time.Duration
	amount      int
}

type alerts []scheduleAlert

type SpyBlindAlerter struct {
	alert alerts
}

func (ba *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	ba.alert = append(ba.alert, scheduleAlert{duration, amount})
}

var dummySpyAlerter = &SpyBlindAlerter{}

func TestGame(t *testing.T) {

	db, err, clean := createTempFile("")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer clean()

	playerStore, e := poker.NewFileSystemPlayerStore(db)
	if e != nil {
		t.Errorf("problem initialising player db file, %v", e)
	}

	for _, name := range []string{"Chris", "pepper"} {
		t.Run("test recording player win", func(t *testing.T) {
			in := strings.NewReader("1\n" + name + " wins\n")
			cli := poker.NewPokerCLI(in, dummyStdOut, poker.NewGame(playerStore, dummySpyAlerter))
			cli.PlayPoker()

			assertPlayerWin(t, playerStore, name)
		})
	}

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\n" + "Chris wins\n")
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewPokerCLI(in, dummyStdOut, poker.NewGame(playerStore, blindAlerter))
		cli.PlayPoker()

		cases := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {

				if len(blindAlerter.alert) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alert)
				}

				alert := blindAlerter.alert[i]

				assertAlert(t, alert, c.expectedScheduleTime, c.expectedAmount)
			})
		}
	})
}

func assertPlayerWin(t *testing.T, store poker.PlayerStore, name string) {
	t.Helper()

	got, err := store.GetScore(name)
	if err != nil {
		t.Fatalf("Failed to record a win, %v", err)
	}
	if got != 1 {
		t.Fatalf("expected a win call for %s, but didn't get any", name)
	}
}

func assertAlert(t *testing.T, alert scheduleAlert, expectedScheduleTime time.Duration,
	expectedAmount int) {
	t.Helper()

	amountGot := alert.amount
	if amountGot != expectedAmount {
		t.Errorf("got amount %d, want %d", amountGot, expectedAmount)
	}

	gotScheduledTime := alert.scheduledAt
	if gotScheduledTime != expectedScheduleTime {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, expectedScheduleTime)
	}
}
