package middleware

import (
	"Server/httpServer"
	"Server/service"
	"net/http"
	"strings"
)

func TokenVerificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			httpServer.Unauthorized(w, "Missing token")
			return
		}

		// Assuming tokens are in the format "Bearer <token>"
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			httpServer.Unauthorized(w, "Invalid token format")
			return
		}

		token = splitToken[1] // The actual token part

		// Validate the token (you would implement your token validation logic here)
		isValid := service.IsTokenValid(token)

		if !isValid {
			httpServer.Unauthorized(w, "Invalid or expired token")
			return
		}

		// If token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
