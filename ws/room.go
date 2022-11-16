package ws

import "github.com/google/uuid"

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

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetTitle() string {
	return room.Title
}
