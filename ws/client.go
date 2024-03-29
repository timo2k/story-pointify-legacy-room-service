package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/story-pointify/room-service/utils"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxPayloadSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type Client struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Title             string    `json:"title"`
	CurrentEstimation string    `json:"current-estimation"`
	IsSpectator       bool      `json:"is-spectator"`
	conn              *websocket.Conn
	wsServer          *WsServer
	send              chan []byte
	rooms             map[*Room]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	return &Client{
		ID:                uuid.New(),
		Name:              name,
		Title:             utils.GenerateRandomDinoName(),
		CurrentEstimation: "0",
		IsSpectator:       false,
		conn:              conn,
		wsServer:          wsServer,
		send:              make(chan []byte, 256),
		rooms:             make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxPayloadSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonPayload, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error %v", err)
			}
			break
		}

		client.handleNewPayload(jsonPayload)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url query param 'name' is missing")
		return
	}

	// Allow All Origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, name[0])

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

func (client *Client) handleNewPayload(jsonPayload []byte) {
	var payload Payload
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	payload.Sender = client

	switch payload.Event {

	case OnSendEstimation:
		client.handleSendEstimation(payload)

		if payload.Target == nil {
			log.Println("Target is null :(")
			return
		}

		roomId := payload.Target.GetId()
		if room := client.wsServer.findRoomById(roomId); room != nil {
			room.broadcast <- &payload
		}

	case OnToggleHideEstimations:

		if payload.Target == nil {
			log.Println("Target is null :(")
			return
		}

		roomId := payload.Target.GetId()
		if room := client.wsServer.findRoomById(roomId); room != nil {
			payload.Target.HasHiddenEstimations = client.handleToggleShowAndHideEstimations(payload)
			room.HasHiddenEstimations = client.handleToggleShowAndHideEstimations(payload)
			room.broadcast <- &payload
		}

	case OnJoinRoom:
		client.handleJoinRoomPayload(payload)

	case OnLeaveRoom:
		client.handleLeaveRoomPayload(payload)
	}
}

func (client *Client) handleSendEstimation(payload Payload) {
	client.CurrentEstimation = payload.Message
}

func (client *Client) handleToggleShowAndHideEstimations(payload Payload) bool {
	switch payload.Message {
	case "hide":
		return true
	case "show":
		return false
	}
	return false
}

func (client *Client) handleJoinRoomPayload(payload Payload) {

	// slice message string into array to get the title and specator flag for join room payload
	titleAndSpectatorFlag := strings.Split(payload.Message, ";")

	if len(titleAndSpectatorFlag) >= 2 {
		roomTitle := titleAndSpectatorFlag[0]
		isSpectator, err := strconv.ParseBool(titleAndSpectatorFlag[1])

		if err != nil {
			log.Println("Cannot cast isSpectator string to boolean")
			client.joinRoom(roomTitle, false, client)
			return
		}

		client.joinRoom(roomTitle, isSpectator, client)
	}

}

func (client *Client) handleLeaveRoomPayload(payload Payload) {
	room := client.wsServer.findRoomById(payload.Message)
	if room == nil {
		return
	}

	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func (client *Client) joinRoom(roomTitle string, isSpectator bool, sender *Client) {
	client.IsSpectator = isSpectator

	room := client.wsServer.findRoomByTitle(roomTitle)
	if room == nil {
		room = client.wsServer.createRoom(roomTitle)
	}

	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client

		client.notifyRoomJoined(room, sender)
	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	payload := Payload{
		Event:  OnRoomJoined,
		Target: room,
		Sender: sender,
	}

	client.send <- payload.encode()
}

func (client *Client) GetId() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}
