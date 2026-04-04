package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func initJWTSecret() {
	pass := os.Getenv("TODO_PASSWORD")
	if pass != "" {
		hash := sha256.Sum256([]byte(pass))
		jwtSecret = []byte(hex.EncodeToString(hash[:]))
	}
}

// signInHandler обрабатывает POST /api/signin
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "authentication not configured"})
		return
	}

	if req.Password != expectedPass {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong password"})
		return
	}

	// Создаём JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": sha256Hash(expectedPass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// auth middleware для проверки аутентификации
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			// Пароль не установлен - аутентификация не требуется
			next(w, r)
			return
		}

		// Проверяем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		// Проверяем хеш пароля в токене
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if hash, ok := claims["hash"].(string); ok {
				if hash != sha256Hash(pass) {
					http.Error(w, "authentication required", http.StatusUnauthorized)
					return
				}
			}
		}

		next(w, r)
	}
}

func sha256Hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
