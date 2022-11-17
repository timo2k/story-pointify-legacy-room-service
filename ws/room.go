package ws

import (
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Payload
}

func NewRoom(title string) *Room {
	return &Room{
		ID:         uuid.New(),
		Title:      title,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Payload),
	}
}

func (room *Room) runRoom() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case payload := <-room.broadcast:
			room.broadcastToClientsInRoom(payload.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastToClientsInRoom(payload []byte) {
	for client := range room.clients {
		client.send <- payload
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	payload := &Payload{
		Event:   OnSendMessage,
		Target:  room,
		Message: fmt.Sprintf("Hallo %s", client.GetName()),
	}

	room.broadcastToClientsInRoom(payload.encode())
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetTitle() string {
	return room.Title
}
