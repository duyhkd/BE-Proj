package handlers

import (
	"Server/httpServer"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type HttpHandler struct {
	redis redis.Client
}

func NewHandler(redis *redis.Client) *HttpHandler {
	return &HttpHandler{
		redis: *redis,
	}
}

func (handler HttpHandler) Ping(c *gin.Context) {
	httpServer.Ok(c, "Server is up!")
}
