package handlers

import (
	"errors"
	"net/http"

	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var (
	ErrPermissionRequired = errors.New("permission required")
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

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		c.JSON(makeHTTPResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	channel, err := handler.channelUsecase.GetChannelByID(c.Request.Context(), channelID)
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	if !handler.channelUsecase.IsMember(c.Request.Context(), channel.ID, userID) {
		c.JSON(makeHTTPResponse(http.StatusForbidden, ErrPermissionRequired.Error(), nil))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(makeHTTPResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	message, err := handler.messageUsecase.CreateNewMessage(c.Request.Context(), &domains.CreateMessageRequest{
		Text:      req.Text,
		File:      req.File,
		ChannelID: channel.ID,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	response := domains.Message{
		ID:        message.ID,
		Text:      message.Text,
		File:      message.File,
		ChannelID: message.ChannelID,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		Member:    message.Member,
	}

	handler.socketService.EmitNewMessage(channelID, &response)

	c.JSON(makeHTTPResponse(http.StatusCreated, "Message created Successfully", message))
}

func (handler *messageHandler) EditMessage(c *gin.Context) {
	messageID := c.Param("messageID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	var req domains.MessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHTTPResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	message, err := handler.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	if message.Member.ID != userID {
		c.JSON(makeHTTPResponse(http.StatusForbidden, ErrPermissionRequired.Error(), nil))
		return
	}

	updatedMessage, err := handler.messageUsecase.UpdateMessageByID(c.Request.Context(), message.ID, *req.Text)
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	response := domains.Message{
		ID:        updatedMessage.ID,
		Text:      updatedMessage.Text,
		File:      updatedMessage.File,
		ChannelID: updatedMessage.ChannelID,
		CreatedAt: updatedMessage.CreatedAt,
		UpdatedAt: updatedMessage.UpdatedAt,
		Member:    message.Member,
	}

	handler.socketService.EmitEditMessage(updatedMessage.ChannelID, &response)

	c.JSON(makeHTTPResponse(http.StatusOK, "Message updated Successfully", updatedMessage))
}

func (handler *messageHandler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageID")
	userID := c.MustGet(middleware.AuthorizationUserIdKey).(string)

	message, err := handler.messageUsecase.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	if message.Member.ID != userID {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, ErrPermissionRequired.Error(), nil))
		return
	}

	if err := handler.messageUsecase.DeleteMessageByID(c.Request.Context(), message.ID); err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	handler.socketService.EmitDeleteMessage(message.ChannelID, message.ID)

	c.JSON(makeHTTPResponse(http.StatusOK, "Message deleted successfully", message))
}
