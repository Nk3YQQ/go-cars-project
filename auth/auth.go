package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("Key")

type contextKey string

const userIDKey contextKey = "userID"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func ParseToken(authToken string) (string, error) {
	// Функция для парсинга Bearer токена

	if len(authToken) <= 7 && authToken[:7] != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return authToken[7:], nil
}

func GenerateToken(userID uint) (string, error) {
	// Функция для генерации jwt токена

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	// Контроллер для аутентификации пользователя

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		claims := &Claims{}

		parsedToken, err := ParseToken(tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		token, err := jwt.ParseWithClaims(parsedToken, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected sign method")
			}

			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
