package clockface

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClockFaceWS struct {
	*websocket.Conn
}

func NewWwebSocket(w http.ResponseWriter, r *http.Request) *ClockFaceWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &ClockFaceWS{conn}
}

/**
Enable it to be able to used for writing/sending msg down to the client's browser
*/
func (ws *ClockFaceWS) Write(p []byte) (n int, err error) {
	err = ws.WriteMessage(1, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
