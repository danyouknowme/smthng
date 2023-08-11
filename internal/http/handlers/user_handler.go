package handlers

import (
	"net/http"
	"time"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/pkg/helpers"
	"github.com/danyouknowme/smthng/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUsecase    usecases.UserUsecase
	sessionUsecase usecases.SessionUsecase
	jwtService     jwt.JWTService
}

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

func NewUserHandler(userUsecase usecases.UserUsecase, sessionUsecase usecases.SessionUsecase, jwtService jwt.JWTService) UserHandler {
	return &userHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		jwtService:     jwtService,
	}
}

func (handler *userHandler) Register(c *gin.Context) {
	var req *domains.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Password = hashedPassword

	err = handler.userUsecase.CreateNewUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})
}

func (handler *userHandler) Login(c *gin.Context) {
	var req *domains.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, err := handler.userUsecase.Authenticate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken, accessTokenPayload, err := handler.jwtService.CreateToken(userID, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken, refreshTokenPayload, err := handler.jwtService.CreateToken(userID, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	session, err := handler.sessionUsecase.CreateNewSession(c.Request.Context(), &domains.SessionMongo{
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIP:     c.ClientIP(),
		ExpiredAt:    refreshTokenPayload.ExpiredAt,
	})

	c.JSON(http.StatusOK, gin.H{
		"session_id":               session.ID,
		"access_token":             accessToken,
		"access_token_expired_at":  accessTokenPayload.ExpiredAt,
		"refresh_token":            refreshToken,
		"refresh_token_expired_at": refreshTokenPayload.ExpiredAt,
	})
}
