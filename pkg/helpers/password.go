package helpers

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if !PasswordCondition(password) {
		return "", fmt.Errorf("password not match with the condition")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func PasswordCondition(password string) bool {
	var (
		upp, low, num, sym bool
		length             uint8
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
			length++
		case unicode.IsLower(char):
			low = true
			length++
		case unicode.IsNumber(char):
			num = true
			length++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			length++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || length < 6 {
		return false
	}

	return true
}
