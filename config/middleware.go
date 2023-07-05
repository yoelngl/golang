package config

import (
	"fmt"
	"net/http"
	"os"
	"restfulapi/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redisServer, _ := ConnectRedis()
		var keyToken = []byte(os.Getenv("JWT_SECRET_KEY"))
		header := strings.TrimSpace(r.Header.Get("Authorization"))

		_, err := redisServer.Get(header).Result()
		if err == nil {
			utils.JsonResponse(w, "error", "Invalid auth token!", http.StatusUnauthorized, nil)
			return
		}

		if header == "" {
			utils.JsonResponse(w, "error", "Missing auth token", http.StatusUnauthorized, nil)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				utils.JsonResponse(w, "error", "unexpected signing method", http.StatusUnauthorized, nil)
				return nil, fmt.Errorf("unexpected signing method")
			}

			return keyToken, nil
		})

		if err != nil {
			utils.JsonResponse(w, "error", "Invalid auth token!", http.StatusUnauthorized, []string{err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			r.Header.Set("User", claims["user"].(string))
			next.ServeHTTP(w, r)
		} else {
			utils.JsonResponse(w, "error", "Invalid auth token!", http.StatusUnauthorized, nil)
			return
		}

	})
}
