package auth

import (
	"fmt"
	"time"

	errs "github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/dgrijalva/jwt-go"
)

// JWTManager is a JSON web token manager
type JWTManager struct {
	secretKey       string
	tokenDuration   time.Duration
	refreshKey      string
	refreshDuration time.Duration
}

// UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	//Role     string `json:"role"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

// NewJWTManager returns a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration, refreshKey string, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration, refreshKey, refreshDuration}
}

//func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
//	return &JWTManager{secretKey, tokenDuration}
//}

// Generate generates and signs a new token for a user
func (manager *JWTManager) Generate(username string) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// GenerateRefresh generates and signs a new token for a user
func (manager *JWTManager) GenerateRefresh(username, email, id string) (string, error) {
	claims := RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.refreshDuration).Unix(),
		},
		Username: username,
		Email:    email,
		ID:       id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.refreshKey))
}

// Verify verifies the access token string and return a user claim if the token is valid
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.Error().Msg("unexpected token signing method")
				return nil, errs.New("Bad Token Signing")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, errs.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// VerifyRefresh verifies the access token string and return a user claim if the token is valid
func (manager *JWTManager) VerifyRefresh(accessToken string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&RefreshClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.refreshKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
