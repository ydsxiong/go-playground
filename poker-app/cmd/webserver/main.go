package main

import (
	"log"
	"net/http"

	"github.com/ydsxiong/go-playground/poker-app/poker"
)

/**
Earlier we explored that the Handler interface is what we need to implement in order to make a server.
Typically we do that by creating a struct and make it implement the interface.
However the use-case for structs is for holding data but currently we have no state,
so it doesn't feel right to be creating one.

HandlerFunc lets us avoid this.
The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.
If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.

type HandlerFunc func(ResponseWriter, *Request)

*/

const dbFileName = "../../game.db.json"

func main() {
	store, closeStore, err := poker.LoadUpFileStore(dbFileName)
	if err != nil {
		log.Fatalf("Problem with loading in file store, %v", err)
	}
	defer closeStore()

	game := poker.NewGame(store, poker.BlindAlerterFunc(poker.Alerter))
	server, e := poker.NewPlayerServer(store, game)
	if e != nil {
		log.Fatalf("Problem with setting up the server, %v", err)
	}

	if http.ListenAndServe(":5000", server) != nil {
		log.Fatalf("Couldn't listen to 5000 port, %v", err)
	}
}
