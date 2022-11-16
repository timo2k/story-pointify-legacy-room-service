package ws

import (
	"encoding/json"
	"log"
)

const (
	SendEstimationAction = "send-estimation"
	JoinRoomAction       = "join-room"
	LeaveRoomAction      = "leave-room"
	UserJoinedAction     = "user-join"
	UserLeftAction       = "user-left"
	RoomJoinedAction     = "room-joined"
)

type Payload struct {
	Action            string  `json:"action"`
	Target            *Room   `json:"target"`
	Sender            *Client `json:"sender"`
	Message           string  `json:"message"`
	CurrentEstimation string  `json:"current-estimation"`
}

func (payload *Payload) encode() []byte {
	json, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	return json
}
