package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	rooms    map[*Room]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		ID:       uuid.New(),
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}
}
