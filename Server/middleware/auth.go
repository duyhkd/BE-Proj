package middleware

import (
	"Server/httpServer"
	"Server/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			httpServer.Unauthorized(c, "Missing token")
			return
		}

		// Assuming tokens are in the format "Bearer <token>"
		splitToken := strings.Split(token, "Bearer ")
		if len(splitToken) != 2 {
			httpServer.Unauthorized(c, "Invalid token format")
			return
		}

		tokenstring := splitToken[1] // The actual token part

		// isValid := service.IsTokenValid(token)

		parsedToken, claims := service.ParseToken(tokenstring)

		if !parsedToken.Valid {
			httpServer.Unauthorized(c, "Invalid or expired token")
		}

		// Store claims in the context for later use
		c.Set("username", claims.Username)
		// ctx := context.WithValue(c.Request.Context(), "username", claims.Username)
		// c.Request.WithContext(ctx)

		// If token is valid, proceed to the next handler
		c.Next()
	}
}
