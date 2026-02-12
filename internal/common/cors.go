package common

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func CORS(next http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:4200",
			"http://localhost:4200", "http://localhost:8081"}),
		handlers.AllowCredentials(),
	)(next)
}
