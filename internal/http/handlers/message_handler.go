package handlers

import (
	"time"

	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type messageHandler struct {
	socketService ws.SocketService
}

type MessageHandler interface {
	CreateNewMessage(c *gin.Context)
}

func NewMessageHandler(socketService ws.SocketService) MessageHandler {
	return &messageHandler{
		socketService: socketService,
	}
}

type messageRequest struct {
	Text string `json:"text"`
}

func (h *messageHandler) CreateNewMessage(c *gin.Context) {
	channelID := c.Param("channelID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	var message messageRequest

	response := domains.Message{
		ID:        "",
		Text:      message.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Member: domains.User{
			ID:           userID,
			Username:     "",
			ProfileImage: "",
			IsOnline:     true,
		},
	}

	h.socketService.EmitNewMessage(channelID, &response)
}
