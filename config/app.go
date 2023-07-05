package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	} else {
		log.Println("Success load .env!")
	}
}

func Databases() {
	var err error
	_, err = ConnectRedis()
	if err != nil {
		log.Fatalf("Failed to connect to database redis: %v", err)
	} else {
		log.Println("Database Redis successfully Connect!")
	}

	_, err = Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database MySQL: %v", err)
	} else {
		log.Println("Database MySQL successfully Connect!")
	}

	return
}
