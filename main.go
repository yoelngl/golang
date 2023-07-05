package main

import (
	"log"
	"net/http"
	"os"
	"restfulapi/config"
	"restfulapi/routes"
	"time"
)

func main() {
	config.LoadEnv()
	config.Databases()
	router := routes.Router()

	port := os.Getenv("PORT")
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Server failed to start:", err)
	}
}
