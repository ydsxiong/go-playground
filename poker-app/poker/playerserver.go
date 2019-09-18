package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/**
We changed the second property of PlayerServer, removing the named property router http.ServeMux
and replaced it with http.Handler; this is called embedding.
Go does not provide the typical, type-driven notion of subclassing, but it does have the ability
to “borrow” pieces of an implementation by embedding types within a struct or interface.

You must be careful with embedding types because you will expose all public methods and fields of
the type you embed. In our case, it is ok because we embedded just the interface that we wanted to expose (http.Handler).
If we had been lazy and embedded http.ServeMux instead (the concrete type) it would still work
but users of PlayerServer would be able to add new routes to our server because Handle(path, handler) would be public.
When embedding types, really think about what impact that has on your public API.
It is a very common mistake to misuse embedding and end up polluting your APIs and exposing the internals of your type.

*/
type PlayerServer struct {
	Store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

type playerServerWS struct {
	*websocket.Conn
}

const htmlTemplatePath = "game.html"

/**

It's quite odd (and inefficient) to be setting up a router as a request comes in and then calling it.
What we ideally want to do is have some kind of NewPlayerServer function which will take our dependencies
and do the one-time setup of creating the router. Each request can then just use that one instance of the router.

func (ps *PlayerServer) Servehttp(w http.ResponseWriter, r *http.Request) {

	router := http.NewServeMux()
	router.HandleFunc("/league", http.HandlerFunc(ps.handleLeague))
	router.HandleFunc("/players/", ps.handlePlayers)

	router.ServeHTTP(w, r)
}
*/

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	ps := new(PlayerServer)
	ps.game = game

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		//http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	ps.template = tmpl
	ps.Store = store

	router := http.NewServeMux()
	router.HandleFunc("/league", http.HandlerFunc(ps.handleLeague))
	router.HandleFunc("/players/", ps.handlePlayers)
	router.HandleFunc("/game", ps.handleGame)
	router.HandleFunc("/ws", ps.webSocket)

	ps.Handler = router

	return ps, nil
}

/**
Remember that every time we get a winner we close the connection, you will need to refresh the page to open the connection again.
*/
func (ps *PlayerServer) handleGame(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("content-type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(ps.Store.GetLeagueTable())
	//if err != nil {
	//http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
	//}
	ps.template.Execute(w, nil)

}

/**
The product owner is thrilled with the command line application but would prefer it if we could bring that functionality to the browser.
She imagines a web page with a text box that allows the user to enter the number of players and when they submit the form the page displays
the blind value and automatically updates it when appropriate. Like the command line application the user can declare the winner
and it'll get saved in the database.
On the face of it, it sounds quite simple but as always we must emphasise taking an iterative approach to writing software.
First of all we will need to serve HTML. So far all of our HTTP endpoints have returned either plaintext or JSON. We could use the same techniques
we know (as they're all ultimately strings) but we can also use the html/template package for a cleaner solution.
We also need to be able to asynchronously send messages to the user saying The blind is now *y* without having to refresh the browser.
We can use WebSockets to facilitate this.
WebSocket is a computer communications protocol, providing full-duplex communication channels over a single TCP connection
Given we are taking on a number of techniques it's even more important we do the smallest amount of useful work possible first and then iterate.
For that reason the first thing we'll do is create a web page with a form for the user to record a winner. Rather than using a plain form,
we will use WebSockets to send that data to our server for it to record.
*/
func (ps *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := NewWwebSocket(w, r)

	playersMsg := ws.WaitForMsg()
	number, err := strconv.Atoi(playersMsg)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid vaue for number (%s) of players %s", playersMsg, err.Error()), http.StatusInternalServerError)
		return
	}
	ps.game.Start(number, ws) //ioutil.Discard)

	winnerMsg := ws.WaitForMsg()
	ps.game.Finish(winnerMsg)
}

func (ps *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(ps.Store.GetLeagueTable())
}

func (ps *PlayerServer) handlePlayers(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		ps.showScore(w, r, player)
	case http.MethodPost:
		ps.processWin(w, r, player)
	}
}

func (ps *PlayerServer) showScore(w http.ResponseWriter, r *http.Request, player string) {
	//fmt.Fprintf(w, "20") // ok
	// fmt.Println(w, "20") // no
	// w.Write([]byte("20")) // ok
	// r.URL.Path[len("/players/"):]

	score, err := ps.Store.GetScore(player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func (ps *PlayerServer) processWin(w http.ResponseWriter, r *http.Request, player string) {
	ps.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
	ps.showScore(w, r, player)
}

func NewWwebSocket(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		//http.Error(w, fmt.Sprintf("problem accepting websocket connection %s", err.Error()), http.StatusInternalServerError)
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (ws *playerServerWS) WaitForMsg() string {
	_, msg, err := ws.Conn.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

func (ws *playerServerWS) Write(p []byte) (n int, err error) {
	err = ws.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
