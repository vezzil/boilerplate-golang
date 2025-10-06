package jwtmanager

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager issues and validates JWT tokens. This is a minimal implementation
// that can be swapped for another provider later.
type Manager struct {
	Secret     []byte
	Issuer     string
	ExpireIn   time.Duration
}

// Claims embeds RegisteredClaims with custom fields if needed.
type Claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

// New creates a new JWT manager instance.
func New(secret, issuer string, expireIn time.Duration) *Manager {
	return &Manager{Secret: []byte(secret), Issuer: issuer, ExpireIn: expireIn}
}

// Sign creates a signed JWT string.
func (m *Manager) Sign(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ExpireIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

// Verify parses and validates a token string.
func (m *Manager) Verify(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
