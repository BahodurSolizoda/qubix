package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
	}
}

func (s *JWTService) GenerateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (int, string, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("invalid signing method")
			}

			return s.secretKey, nil
		},
	)

	if err != nil || !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	userIDValue, ok := claims["user_id"].(float64)
	if !ok || userIDValue <= 0 {
		return 0, "", errors.New("invalid user ID")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return 0, "", errors.New("invalid user role")
	}

	return int(userIDValue), role, nil
}
