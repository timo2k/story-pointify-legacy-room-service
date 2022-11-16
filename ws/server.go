package ws

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewWebSocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case client := <-server.broadcast:
			server.broadcastToClients(client)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.notifyClientJoined(client)
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		server.notifyClientLeft(client)
	}
}

func (server *WsServer) notifyClientJoined(client *Client) {
	payload := &Payload{
		Event:  OnUserJoined,
		Sender: client,
	}

	server.broadcastToClients(payload.encode())
}

func (server *WsServer) notifyClientLeft(client *Client) {
	payload := &Payload{
		Event:  OnUserLeft,
		Sender: client,
	}

	server.broadcastToClients(payload.encode())
}

func (server *WsServer) broadcastToClients(payload []byte) {
	for client := range server.clients {
		client.send <- payload
	}
}

func (server *WsServer) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientById(id string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.GetId() == id {
			foundClient = client
			break
		}
	}

	return foundClient
}

func (server *WsServer) findRoomByTitle(title string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetTitle() == title {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByName(name string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.GetName() == name {
			foundClient = client
			break
		}
	}

	return foundClient
}

func (server *WsServer) createRoom(title string) *Room {

	return nil
}
