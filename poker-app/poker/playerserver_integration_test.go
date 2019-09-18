package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/ydsxiong/go-playground/poker-app/poker"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, err, cleanDatabase := createTempFile("")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer cleanDatabase()
	store, e := poker.NewFileSystemPlayerStore(database)
	if e != nil {
		t.Fatalf("Problem with opening file database, %v", e)
	}
	server := mustMakePlayerServer(t, store, dummyGame)
	player := "pepper"

	server.ServeHTTP(httptest.NewRecorder(), createRequestFor(http.MethodPost, player))
	server.ServeHTTP(httptest.NewRecorder(), createRequestFor(http.MethodPost, player))
	server.ServeHTTP(httptest.NewRecorder(), createRequestFor(http.MethodPost, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, createRequestFor(http.MethodGet, player))
		if "3" != response.Body.String() {
			t.Errorf("got %s, but wanted %s", response.Body.String(), "3")
		}
		if response.Code != http.StatusOK {
			t.Errorf("wrong status code, got %d, but wanted %d", response.Code, http.StatusOK)
		}
	})

	t.Run("get league", func(t *testing.T) {
		res, _ := createandservereqres(server)
		got := assertresponsecode(t, res)

		wanted := poker.League{
			{"pepper", 3},
		}
		assertleague(t, got, wanted)
	})
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}
