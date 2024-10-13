package handlers

import (
	"Server/httpServer"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type HttpHandler struct {
	redis redis.Client
}

const (
	topPingKey   = "top"
	pingCountKey = "user_ping_count"
)

var mutex sync.Mutex

func NewHandler(redisClient redis.Client) *HttpHandler {
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	return &HttpHandler{
		redis: redisClient,
	}
}

func (handler HttpHandler) Ping(c *gin.Context) {
	username := c.Query("username")

	mutex.Lock()
	defer mutex.Unlock()

	// track unique users with HyperLogLog
	handler.redis.PFAdd(pingCountKey, username)

	// Top Ping tracking
	handler.redis.ZIncrBy(topPingKey, 1, username)

	// sleep 5s
	time.Sleep(5 * time.Second)
	httpServer.Ok(c, "Server is up!")
}

func (handler HttpHandler) TopPing(c *gin.Context) {
	topUsers, _ := handler.redis.ZRevRangeWithScores(topPingKey, 0, 9).Result()
	result := make([]gin.H, len(topUsers))
	for i, user := range topUsers {
		result[i] = gin.H{"username": user.Member, "count": user.Score}
	}

	// Need to fix response
	c.JSON(http.StatusOK, result)
}

func (handler HttpHandler) PingCount(c *gin.Context) {
	count, err := handler.redis.PFCount(pingCountKey).Result()
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	}
	httpServer.Ok(c, fmt.Sprintf("Ping Count: %v", count))
}
