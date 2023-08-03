package handlers

import (
	"net/http"

	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type messageHandler struct {
	socketService  ws.SocketService
	channelUsecase usecases.ChannelUsecase
	messageUsecase usecases.MessageUsecase
}

type MessageHandler interface {
	CreateNewMessage(c *gin.Context)
}

func NewMessageHandler(
	socketService ws.SocketService,
	channelUsecase usecases.ChannelUsecase,
	messageUsecase usecases.MessageUsecase,
) MessageHandler {
	return &messageHandler{
		socketService:  socketService,
		channelUsecase: channelUsecase,
		messageUsecase: messageUsecase,
	}
}

type messageRequest struct {
	Text string `json:"text"`
}

func (handler *messageHandler) CreateNewMessage(c *gin.Context) {
	channelID := c.Param("channelID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	var req messageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	channel, err := handler.channelUsecase.GetChannelByID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !handler.channelUsecase.IsMember(c.Request.Context(), channel.ID, userID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not a member of this channel",
		})
		return
	}

	message, err := handler.messageUsecase.CreateNewMessage(c.Request.Context(), &domains.CreateMessageRequest{
		Text:      req.Text,
		ChannelID: channel.ID,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := domains.Message{
		ID:        message.ID,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		Member: domains.User{
			ID:           message.Member.ID,
			Username:     message.Member.Username,
			ProfileImage: message.Member.ProfileImage,
			IsOnline:     message.Member.IsOnline,
		},
	}

	handler.socketService.EmitNewMessage(channelID, &response)

	c.JSON(http.StatusOK, gin.H{
		"message": "message created",
	})
}
