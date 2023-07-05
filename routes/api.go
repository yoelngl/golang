package routes

import (
	"log"
	"net/http"
	"restfulapi/config"
	"restfulapi/controllers"

	"github.com/gorilla/mux"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("New Request: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/api/", controllers.HomeIndex)

	r.HandleFunc("/v1/api/auth/register", controllers.Register)
	r.HandleFunc("/v1/api/auth/login", controllers.Login)
	r.Handle("/v1/api/auth/logout", config.AuthMiddleware(http.HandlerFunc(controllers.Logout)))

	r.Handle("/v1/api/users", config.AuthMiddleware(http.HandlerFunc(controllers.HomeIndex)))

	r.Use(LoggingMiddleware)

	return r

}
