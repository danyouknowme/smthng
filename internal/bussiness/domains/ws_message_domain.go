package domains

import (
	"encoding/json"

	"github.com/danyouknowme/smthng/pkg/logger"
)

type WebsocketMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type ReceivedMessage struct {
	Action  string `json:"action"`
	RoomID  string `json:"room_id"`
	Message *any   `json:"message"`
}

func (message *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		logger.Error(err)
	}

	return encoding
}
