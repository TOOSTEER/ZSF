package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload used across the app.
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"` // "host" | "member"
	TeamID int    `json:"team_id"`
	jwt.RegisteredClaims
}

// Config holds JWT settings loaded from the environment.
type Config struct {
	Secret      string
	ExpiryHours int
}

// GenerateToken creates a signed JWT for the given user.
func GenerateToken(cfg Config, userID int, role string, teamID int) (string, error) {
	expiry := time.Duration(cfg.ExpiryHours) * time.Hour
	if cfg.ExpiryHours <= 0 {
		expiry = 24 * time.Hour // TODO: уточнить дефолтный срок жизни токена
	}

	claims := Claims{
		UserID: userID,
		Role:   role,
		TeamID: teamID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken validates the token string and returns its claims.
func ParseToken(cfg Config, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
