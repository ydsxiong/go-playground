package poker_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ydsxiong/go-playground/poker-app/poker"
)

var (
	dummyGame = &GameSpy{}
	tenMS     = 10 * time.Millisecond
)

func TestGETPlayers(t *testing.T) {

	database, err, cleanDatabase := createTempFile(`[
		{"Name": "pepper", "Wins": 20},
		{"Name": "floyd", "Wins": 10}]`)
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer cleanDatabase()
	store, e := poker.NewFileSystemPlayerStore(database)
	if e != nil {
		t.Fatalf("Problem with opening file database, %v", e)
	}
	server := mustMakePlayerServer(t, store, dummyGame)

	testcases := []struct {
		name     string
		player   string
		req      *http.Request
		w        *httptest.ResponseRecorder
		expected string
		status   int
	}{
		{
			"test1",
			"pepper",
			createRequestFor(http.MethodGet, "pepper"),
			httptest.NewRecorder(),
			"20",
			http.StatusOK,
		},
		{
			"test2",
			"floyd",
			createRequestFor(http.MethodGet, "floyd"),
			httptest.NewRecorder(),
			"10",
			http.StatusOK,
		},
		{
			"test3",
			"foo",
			createRequestFor(http.MethodGet, "foo"),
			httptest.NewRecorder(),
			"0",
			http.StatusNotFound,
		},
		{
			"test4",
			"pepper",
			createRequestFor(http.MethodPost, "pepper"),
			httptest.NewRecorder(),
			"21", // original 20 plus this newly registered one
			http.StatusAccepted,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server.ServeHTTP(tc.w, tc.req)

			got := tc.w.Body.String()
			if got != tc.expected {
				t.Errorf("got %s, but wanted %s", got, tc.expected)
			}
			if tc.w.Code != tc.status {
				t.Errorf("wrong status code, got %d, but wanted %d", tc.w.Code, tc.status)
			}
		})
	}

	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, store, dummyGame)
		res, _ := createandservegamereqres(server)

		if res.Code != http.StatusOK {
			t.Errorf("wrong status code, got %d, but wanted %d", res.Code, http.StatusOK)
		}
	})

	// t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
	// 	plaerserver := mustMakePlayerServer(t, store, dummyGame)
	// 	winner := "Ruth"
	// 	server := httptest.NewServer(plaerserver)
	// 	defer server.Close()

	// 	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	// 	ws := mustDialWS(t, wsURL)
	// 	defer ws.Close()
	// 	writeWSMessage(t, ws, "3")
	// 	writeWSMessage(t, ws, winner)

	// 	time.Sleep(10 * time.Millisecond)
	// 	got, _ := store.GetScore(winner)

	// 	if got != 1 {
	// 		t.Errorf("expected a win score for player %s, but got %d", winner, got)
	// 	}
	// })

	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}

		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, store, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)

		within(t, tenMS, func() { assertMessageReceived(t, ws, wantedBlindAlert) })

	})

}

func within(t *testing.T, waitFor time.Duration, assert func()) {
	t.Helper()

	doneChan := make(chan struct{})

	go func() {
		assert()
		doneChan <- struct{}{}
	}()

	select {
	case <-time.After(waitFor):
		t.Error("timed out")
	case <-doneChan:
	}
}

func assertMessageReceived(t *testing.T, ws *websocket.Conn, wantedBlindAlert string) {
	t.Helper()
	_, gotBlindAlert, _ := ws.ReadMessage()

	if string(gotBlindAlert) != wantedBlindAlert {
		t.Errorf("got blind alert %q, want %q", string(gotBlindAlert), wantedBlindAlert)
	}
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func createRequestFor(method, name string) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func createandservereqres(server *poker.PlayerServer) (*httptest.ResponseRecorder, *http.Request) {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	return res, req
}

func createandservegamereqres(server *poker.PlayerServer) (*httptest.ResponseRecorder, *http.Request) {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	return res, req
}

func assertresponsecode(t *testing.T, res *httptest.ResponseRecorder) poker.League {
	t.Helper()
	var got poker.League

	err := json.NewDecoder(res.Body).Decode(&got)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", res.Body, err)
	}

	if res.Code != http.StatusOK {
		t.Errorf("wrong status code, got %d, but wanted %d", res.Code, http.StatusOK)
	}

	if res.Result().Header.Get("content-type") != "application/json" {
		t.Errorf("response did not have content-type of application/json, got %v", res.Result().Header)
	}

	return got
}

func assertleague(t *testing.T, got, wanted poker.League) {
	t.Helper()
	if !reflect.DeepEqual(got, wanted) { //if !playerstore.AreCollectionsEqual(got, wanted) {
		t.Errorf("got %v want %v", got, wanted)
	}
}

func areSameCollections(a, b poker.League) bool {
	common := intersection(a, b)
	expected := len(b)
	return len(a) == expected && len(common) == expected
}

func intersection(a, b poker.League) (c poker.League) {
	m := make(map[poker.Player]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}
