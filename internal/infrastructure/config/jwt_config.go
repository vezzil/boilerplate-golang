package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"boilerplate-golang/internal/infrastructure/jwtmanager"
)

// JWTConfig holds the JWT configuration
var (
	JWT            *jwtmanager.Manager
	RefreshJWT     *jwtmanager.Manager
	RefreshSecrets = make(map[uint]string) // In production, use Redis or database
)

// InitJWT initializes the JWT manager with configuration
func InitJWT() {
	cfg := Get()

	// Generate a secure random secret if not set
	if cfg.JWT.Secret == "" {
		secret, err := generateRandomKey(64) // 64 bytes = 512 bits
		if err != nil {
			panic(fmt.Sprintf("Failed to generate JWT secret: %v", err))
		}
		cfg.JWT.Secret = secret
	}

	// Initialize access token manager
	JWT = jwtmanager.New(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.ExpireIn,
	)

	// Initialize refresh token manager with longer expiration (7 days)
	RefreshJWT = jwtmanager.New(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		24*time.Hour*7, // 7 days
	)
}

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GenerateTokenPair generates a new access token and refresh token for a user
func GenerateTokenPair(userID uint) (*TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := generateToken(userID, JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := generateToken(userID, RefreshJWT)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash the refresh token for secure storage
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	// In production, store this in Redis or database with user ID and expiration
	RefreshSecrets[userID] = string(hashedToken)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// VerifyRefreshToken verifies a refresh token and returns a new token pair
func VerifyRefreshToken(userID uint, refreshToken string) (*TokenPair, error) {
	// In production, verify the refresh token against the stored hash in Redis/database
	storedHash, exists := RefreshSecrets[userID]
	if !exists {
		return nil, fmt.Errorf("no refresh token found for user")
	}

	// Verify the token hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refreshToken)); err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Generate new token pair
	return GenerateTokenPair(userID)
}

// InvalidateRefreshToken removes a user's refresh token
func InvalidateRefreshToken(userID uint) {
	delete(RefreshSecrets, userID)
}

// generateToken is a helper function to generate a JWT token
func generateToken(userID uint, manager *jwtmanager.Manager) (string, time.Time, error) {
	tokenString, err := manager.Sign(userID)
	if err != nil {
		return "", time.Time{}, err
	}

	// Parse the token to get expiration time
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &jwtmanager.Claims{})
	if err != nil {
		return "", time.Time{}, err
	}

	claims, ok := token.Claims.(*jwtmanager.Claims)
	if !ok {
		return "", time.Time{}, jwt.ErrInvalidKey
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt.Time, nil
}

// generateRandomKey generates a secure random key of the specified length
func generateRandomKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}
