package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"restfulapi/config"
	"restfulapi/model"
	"restfulapi/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != "POST" {
		utils.JsonResponse(w, "error", "Only accept POST requests", http.StatusMethodNotAllowed, nil)
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JsonResponse(w, "error", "Failed to decode the request body", http.StatusBadRequest, []string{err.Error()})
		return
	}

	user := model.RegisterRequest{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
		Image:    "default.jpg",
	}

	err = utils.ValidateField(user)
	if err != nil {
		utils.JsonResponse(w, "error", err.Error(), http.StatusInternalServerError, []string{err.Error()})
		return
	}

	register, err := model.RegisterStore(user)
	if err != nil {
		utils.JsonResponse(w, "error", register, http.StatusBadRequest, []string{err.Error()})
		return
	}

	utils.JsonResponse(w, "success", register, http.StatusCreated, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		utils.JsonResponse(w, "error", "Only accept POST requests", http.StatusMethodNotAllowed, nil)
		return
	}
	var req model.LoginRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JsonResponse(w, "error", "Failed to decode the request body", http.StatusBadRequest, []string{err.Error()})
		return
	}

	user := model.LoginRequest{
		Email:    strings.TrimSpace(req.Email),
		Password: strings.TrimSpace(req.Password),
	}

	err = utils.ValidateField(user)
	if err != nil {
		utils.JsonResponse(w, "error", err.Error(), http.StatusInternalServerError, []string{err.Error()})
		return
	}

	checkAccount, err := model.CheckAccount(user)
	if err != nil {
		utils.JsonResponse(w, "error", checkAccount, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	if checkAccount == "success" {
		tokenKey := os.Getenv("JWT_SECRET_KEY")
		token := jwt.New(jwt.SigningMethodHS256)

		claims := token.Claims.(jwt.MapClaims)

		claims["authorized"] = true
		claims["user"] = user.Email
		claims["expires"] = time.Now().Add(time.Hour * 2).Unix()

		tokenString, err := token.SignedString([]byte(tokenKey))
		if err != nil {
			utils.JsonResponse(w, "error", "Something error when create a token!", http.StatusInternalServerError, []string{err.Error()})
		}

		utils.JsonResponseWithData(w, "success", "Successfully Login!", tokenString, http.StatusOK, nil)
	} else {
		utils.JsonResponse(w, "error", "Something error in server!", http.StatusInternalServerError, []string{err.Error()})
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.JsonResponse(w, "error", "Only accept POST requests", http.StatusMethodNotAllowed, nil)
		return
	}
	redisServer, _ := config.ConnectRedis()
	tokenString := r.Header.Get("Authorization")

	err := redisServer.Set(tokenString, true, 0).Err()
	if err != nil {
		utils.JsonResponse(w, "error", "Error blacklisting token", http.StatusUnauthorized, []string{err.Error()})
		return
	}

	utils.JsonResponse(w, "success", "Successfully Logout!", http.StatusOK, nil)
}
