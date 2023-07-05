package utils

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	Type    string      `json:"type"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JsonResponse(w http.ResponseWriter, Type string, Message string, Status int, Error []string) {
	log := logrus.New()

	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("errors_logs.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	if Status >= 400 {
		log.WithFields(logrus.Fields{
			"status":             Status,
			logrus.FieldKeyLevel: Type,
			"error":              Error,
		}).Error(Message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Status)
	response := &Response{Type: Type, Status: Status, Message: Message}
	json.NewEncoder(w).Encode(response)
}

func JsonResponseWithData(w http.ResponseWriter, Type string, Message string, Data interface{}, Status int, Error []string) {
	log := logrus.New()

	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("errors_logs.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	if Status >= 400 {
		log.WithFields(logrus.Fields{
			"status":             Status,
			logrus.FieldKeyLevel: Type,
			"error":              Error,
		}).Error(Message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(Status)
	response := &ResponseWithData{Type: Type, Status: Status, Message: Message, Data: Data}
	json.NewEncoder(w).Encode(response)
}
