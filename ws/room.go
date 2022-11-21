package ws

import (
	"log"

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
	if len(room.clients) != 0 {
		room.notifyClientJoinedRoom(client)
	}

	room.clients[client] = true

	room.listOnlineClients(client)
	log.Printf("User: %s successfully joined Room: %s", client.GetName(), room.GetTitle())
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
		room.notifyClientLeftRoom(client)
		log.Printf("User: %s successfully leaved Room", client.GetName())
	}
}

func (room *Room) broadcastToClientsInRoom(payload []byte) {
	for client := range room.clients {
		client.send <- payload
	}
}

func (room *Room) notifyClientJoinedRoom(client *Client) {
	payload := &Payload{
		Event:  OnUserRoomJoined,
		Sender: client,
	}

	room.broadcastToClientsInRoom(payload.encode())
}

func (room *Room) notifyClientLeftRoom(client *Client) {
	payload := &Payload{
		Event:  OnUserRoomLeft,
		Sender: client,
	}

	room.broadcastToClientsInRoom(payload.encode())
}

func (room *Room) listOnlineClients(client *Client) {
	for existingClient := range room.clients {
		payload := &Payload{
			Event:  OnListOnlineClients,
			Sender: existingClient,
		}
		client.send <- payload.encode()
	}
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetTitle() string {
	return room.Title
}
