package api

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func GenerateJWT(secretPassword string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretPassword))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretPassword := os.Getenv("TODO_PASSWORD")
		tokenString, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		if tokenString.Value == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretPassword), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
