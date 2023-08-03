package ws

import (
	"encoding/json"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"github.com/danyouknowme/smthng/pkg/logger"
)

type socketService struct {
	Hub               *Hub
	ChannelRepository repositories.ChannelRepository
}

type SocketService interface {
	EmitNewMessage(room string, message *domains.Message)
}

func NewSocketService(hub *Hub, channelRepository repositories.ChannelRepository) SocketService {
	return &socketService{
		Hub:               hub,
		ChannelRepository: channelRepository,
	}
}

func (s *socketService) EmitNewMessage(roomID string, message *domains.Message) {
	data, err := json.Marshal(domains.WebsocketMessage{
		Action: NewMessageAction,
		Data:   message,
	})

	if err != nil {
		logger.Infof("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, roomID)
}
