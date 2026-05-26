package utils

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// CorsMiddleware returns a new mux.MiddlewareFunc that applies CORS middleware.
func CorsMiddleware() mux.MiddlewareFunc {

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},                                       // Allow requests from any origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow browser preflight
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}

	return cors.New(corsOptions).Handler
}
