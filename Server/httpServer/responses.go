package httpServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MethodNotAllowed(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Method not allowed")
}

func BadRequest(c *gin.Context, message string) {
	c.String(http.StatusBadRequest, message)
}

func Ok(c *gin.Context, message string) {
	c.String(http.StatusOK, message)
}

func StatusInternalServerError(c *gin.Context, message string) {
	c.String(http.StatusInternalServerError, message)
}

func Unauthorized(c *gin.Context, message string) {
	c.String(http.StatusUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	c.String(http.StatusNotFound, message)
}
