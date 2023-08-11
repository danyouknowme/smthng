package apperrors

import "errors"

var (
	ErrHeaderNotProvided   = errors.New("authorization header is not provided")
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")
	ErrUnsupportedAuthType = errors.New("unsupported authorization type")
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

var (
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrPermissionRequired        = errors.New("permission required")
)
