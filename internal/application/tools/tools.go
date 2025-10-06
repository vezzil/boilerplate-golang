package tools

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewUuid generates a new UUID v4 string
func NewUuid() string {
	id := uuid.NewV4()
	return id.String()
}

// IsValidEmail checks if the email has a valid format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// CreateCode generates a random 4-digit code
func CreateCode() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
