package ws

import (
	"encoding/json"
	"log"
)

const (
	OnSendEstimation        = "send-estimation"
	OnResetEstimations      = "reset-estimations"
	OnHideEstimations       = "hide-estimations"
	OnRevealEstimations     = "reveal-estimations"
	OnSendMessage           = "send-message"
	OnJoinRoom              = "join-room"
	OnLeaveRoom             = "leave-room"
	OnUserJoinedServer      = "user-joined-server"
	OnUserLeftServer        = "user-left-server"
	OnUserRoomJoined        = "user-room-joined"
	OnUserRoomLeft          = "user-room-left"
	OnListOnlineClients     = "list-online-clients"
	OnRoomJoined            = "room-joined"
	OnToggleHideEstimations = "toggle-hide-estimations"
)

type Payload struct {
	Event   string  `json:"event"`
	Target  *Room   `json:"target"`
	Sender  *Client `json:"sender"`
	Message string  `json:"message"`
}

func (payload *Payload) encode() []byte {
	json, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	return json
}
