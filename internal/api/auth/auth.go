package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url, omitempty"`
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := ErrorResponse{
				Code:        http.StatusUnauthorized,
				Message:     "Authorization token is required",
				RedirectURL: "/login.html",
			}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		if tokenString.Value == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := ErrorResponse{
				Code:        http.StatusUnauthorized,
				Message:     "Authorization token is required",
				RedirectURL: "/login.html",
			}
			json.NewEncoder(w).Encode(errorResponse)
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
