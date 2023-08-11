package handlers

import (
	"net/http"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type channelHandler struct {
	channelUsecase usecases.ChannelUsecase
}

type ChannelHandler interface {
	CreateNewChannel(c *gin.Context)
}

func NewChannelHandler(channelUsecase usecases.ChannelUsecase) ChannelHandler {
	return &channelHandler{
		channelUsecase: channelUsecase,
	}
}

func (handler *channelHandler) CreateNewChannel(c *gin.Context) {
	var req *domains.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHTTPResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	ownerID := c.MustGet(middleware.AuthorizationUserIdKey).(string)
	req.Members = append(req.Members, ownerID)

	err := handler.channelUsecase.CreateNewChannel(c.Request.Context(), req)
	if err != nil {
		c.JSON(makeHTTPResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(makeHTTPResponse(http.StatusCreated, "Channel created successfully", nil))
}
