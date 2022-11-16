package ws

import (
	"encoding/json"
	"log"
)

const (
	OnSendEstimation    = "send-estimation"
	OnResetEstimations  = "reset-estimations"
	OnHideEstimations   = "hide-estimations"
	OnRevealEstimations = "reveal-estimations"
	OnJoinRoom          = "join-room"
	OnLeaveRoom         = "leave-room"
	OnUserJoined        = "user-join"
	OnUserLeft          = "user-left"
	OnRoomJoined        = "room-joined"
)

type Payload struct {
	Event             string  `json:"event"`
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
