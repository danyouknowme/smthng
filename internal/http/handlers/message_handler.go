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
	CreateMessage(c *gin.Context)
	EditMessage(c *gin.Context)
	DeleteMessage(c *gin.Context)
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

func (handler *messageHandler) CreateMessage(c *gin.Context) {
	channelID := c.Param("channelID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	var req domains.MessageRequest

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
		Member:    message.Member,
	}

	handler.socketService.EmitNewMessage(channelID, &response)

	c.JSON(http.StatusOK, gin.H{
		"message": "message created",
		"data":    message,
	})
}

func (handler *messageHandler) EditMessage(c *gin.Context) {
	messageID := c.Param("messageID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	var req domains.MessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	message, err := handler.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if message.Member.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this message",
		})
		return
	}

	updatedMessage, err := handler.messageUsecase.UpdateMessageByID(c.Request.Context(), message.ID, req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := domains.Message{
		ID:        updatedMessage.ID,
		Text:      updatedMessage.Text,
		CreatedAt: updatedMessage.CreatedAt,
		UpdatedAt: updatedMessage.UpdatedAt,
		Member:    message.Member,
	}

	handler.socketService.EmitEditMessage(updatedMessage.ChannelID, &response)

	c.JSON(http.StatusOK, gin.H{
		"message": "message updated",
		"data":    updatedMessage,
	})
}

func (handler *messageHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	message, err := handler.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if message.Member.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this message",
		})
		return
	}

	if err := handler.messageUsecase.DeleteMessageByID(c.Request.Context(), message.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	handler.socketService.EmitDeleteMessage(message.ChannelID, message.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "message deleted",
		"data": gin.H{
			"message_id": message.ID,
		},
	})
}
