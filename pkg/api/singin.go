package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	expected := os.Getenv("TODO_PASSWORD")
	if expected == "" || req.Password != expected {
		writeError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": fmt.Sprintf("%x", sha256.Sum256([]byte(expected))),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Token generation error")
		return
	}

	writeJSON(w, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			return []byte("secret_key"), nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "Authentication failed")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(pass)))
		if claims["hash"] != hash {
			writeError(w, http.StatusUnauthorized, "Token not valid anymore")
			return
		}

		next(w, r)
	}
}
