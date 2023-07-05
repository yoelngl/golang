package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	host     string
	port     int
	username string
	password string
	database string
}

var client *redis.Client

func Connect() (*sql.DB, error) {
	var databaseConfig Database

	databaseConfig.host = os.Getenv("DATABASE_HOST")
	databaseConfig.port, _ = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	databaseConfig.username = os.Getenv("DATABASE_USER")
	databaseConfig.password = os.Getenv("DATABASE_PASSWORD")
	databaseConfig.database = os.Getenv("DATABASE_NAME")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", databaseConfig.username, databaseConfig.password, databaseConfig.host, databaseConfig.port, databaseConfig.database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
		return nil, err
	}

	return client, nil

}
