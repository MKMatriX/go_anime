package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go_anime/internal/shared/models"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	claims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return nil, nil, err
	}

	refreshAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
		},
	})
	signedRefreshToken, err := refreshAccessToken.SignedString(secretKey)
	if err != nil {
		return nil, nil, err
	}

	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWT(token string) (*CustomJWTClaims, error) {
	parsedJWT, err := jwt.ParseWithClaims(token, &CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	claims, ok := parsedJWT.Claims.(*CustomJWTClaims)
	if !ok {
		return nil, errors.New("Unknown claims type, cannon proceed")
	}

	return claims, nil
}

func IsClaimsExpired(claims *CustomJWTClaims) bool {
	return claims.ExpiresAt.Before(time.Now())
}

func IsTokenExpired(token string) bool {
	claims, err := ParseJWT(token)
	if err != nil {
		slog.Error("Parse claims error: " + err.Error())
	}

	fmt.Printf("Claims: %v+", claims)

	return claims.ExpiresAt.Before(time.Now())
}

func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
