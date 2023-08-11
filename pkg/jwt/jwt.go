package jwt

import (
	"errors"
	"time"

	"github.com/danyouknowme/smthng/pkg/apperrors"
	"github.com/golang-jwt/jwt"
)

type jwtService struct {
	secretKey string
}

type JWTService interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{secretKey}
}

func (maker *jwtService) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *jwtService) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, apperrors.ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, apperrors.ErrExpiredToken) {
			return nil, apperrors.ErrExpiredToken
		}
		return nil, apperrors.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}

	return payload, nil
}
