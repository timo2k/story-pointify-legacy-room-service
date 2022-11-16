package main

import (
	"net/http"

	"github.com/story-pointify/room-service/ws"
)

func main() {
	wsServer := ws.NewWebSocketServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(wsServer, w, r)
	})

	http.ListenAndServe(":1337", nil)
}
