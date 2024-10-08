package handlers

import (
	"Server/httpServer"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	httpServer.Ok(c, "Server is up!")
}
