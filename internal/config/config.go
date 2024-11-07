package config

import (
	"github.com/rs/cors"
	"net/http"
)

func Cors(mux http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization, X-Requested-With,Access-Control-Request-Method, Access-Control-Request-Headers, Access-Control-Allow-Origin"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	return handler
}
