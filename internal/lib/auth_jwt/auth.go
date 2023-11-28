package auth_jwt

import (
	"crypto/sha256"
	"fmt"
	"github.com/Verce11o/yata-auth/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type JWTService struct {
	config config.JWTConfig
}

func NewJWTService(JWTConfig config.JWTConfig) *JWTService {
	return &JWTService{config: JWTConfig}
}

func (j *JWTService) GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.config.TokenTTLHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userId,
	})

	return token.SignedString([]byte(j.config.Secret))
}

func (j *JWTService) ParseToken(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	return claims.UserID, nil
}

func (j *JWTService) GenerateHashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(j.config.Salt)))
}
