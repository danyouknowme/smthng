package helpers

import (
	"fmt"
	"math/rand"
)

func GenerateTag() string {
	randomNumber := rand.Intn(10000)
	return fmt.Sprintf("%04d", randomNumber)
}
