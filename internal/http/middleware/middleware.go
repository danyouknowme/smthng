package middleware

import (
	"net/http"
	"strings"

	"github.com/danyouknowme/smthng/pkg/apperrors"
	"github.com/danyouknowme/smthng/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	AuthorizationUserIdKey  = "userId"
)

func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(makeMiddlewareResponse(apperrors.ErrHeaderNotProvided.Error()))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(makeMiddlewareResponse(apperrors.ErrInvalidHeaderFormat.Error()))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(makeMiddlewareResponse(apperrors.ErrUnsupportedAuthType.Error()))
			return
		}

		accessToken := fields[1]
		payload, err := jwtService.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(makeMiddlewareResponse(err.Error()))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Set(AuthorizationUserIdKey, payload.UserID)
		ctx.Next()
	}
}

type middlewareResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Result     any    `json:"result"`
}

func makeMiddlewareResponse(message string) (int, middlewareResponse) {
	return http.StatusUnauthorized, middlewareResponse{
		Status:     http.StatusText(http.StatusUnauthorized),
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Result:     nil,
	}
}
