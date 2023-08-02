package handlers

import (
	"net/http"

	"github.com/danyouknowme/smthng/internal/bussiness/domains"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/pkg/helpers"
	"github.com/danyouknowme/smthng/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userUsecase usecases.UserUsecase
	jwtService  jwt.JWTService
}

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

func NewUserHandler(userUsecase usecases.UserUsecase, jwtService jwt.JWTService) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
		jwtService:  jwtService,
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

	token, err := handler.jwtService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"CSRF_TOKEN": token,
	})
}
