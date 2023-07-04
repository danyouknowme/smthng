package v1

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
	Room    string `json:"room"`
	Message *any   `json:"message"`
}

func (message *WebsocketMessage) Encode() []byte {
	encoding, err := json.Marshal(message)
	if err != nil {
		logger.Error(err)
	}

	return encoding
}
