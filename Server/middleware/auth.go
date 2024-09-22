package middleware

import (
	"Server/httpServer"
	"Server/service"
	"context"
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

		tokenstring := splitToken[1] // The actual token part

		// isValid := service.IsTokenValid(token)

		parsedToken, claims := service.ParseToken(tokenstring)

		if !parsedToken.Valid {
			httpServer.Unauthorized(w, "Invalid or expired token")
		}

		// Store claims in the context for later use
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		r = r.WithContext(ctx)

		// If token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
