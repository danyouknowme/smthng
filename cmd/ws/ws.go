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

func (s *socketService) EmitEditMessage(roomID string, message *domains.Message) {
	data, err := json.Marshal(domains.WebsocketMessage{
		Action: EditMessageAction,
		Data:   message,
	})

	if err != nil {
		logger.Infof("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, roomID)
}

func (s *socketService) EmitDeleteMessage(roomID string, messageID string) {
	data, err := json.Marshal(domains.WebsocketMessage{
		Action: DeleteMessageAction,
		Data:   messageID,
	})

	if err != nil {
		logger.Infof("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, roomID)
}
