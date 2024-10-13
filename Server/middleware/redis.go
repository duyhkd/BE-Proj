package middleware

import (
	"Server/httpServer"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	redisrate "github.com/go-redis/redis_rate"
)

type RedisMiddleware struct {
	redisClient redis.Client
	limiter     RedisRateLimiter
}
type RedisRateLimiter struct {
	*redisrate.Limiter
}

const rateLimitePrefix = "limiter_%v"

func NewRedisMiddleware(
	redisClient *redis.Client,
) RedisMiddleware {
	limiter := RedisRateLimiter{redisrate.NewLimiter(redisClient)}
	return RedisMiddleware{
		redisClient: *redisClient,
		limiter:     limiter,
	}
}

// Verify Redis Cache
func (m RedisMiddleware) VerifyRedisCache() gin.HandlerFunc {
	return func(c *gin.Context) {

		//Get current URL
		endpoint := c.Request.URL

		//Keys are string typed
		cachedKey := endpoint.String()

		//Get Cached keys
		val, err := m.redisClient.Get(cachedKey).Bytes()

		//If error is nil, it means that redis cache couldn't find the key, hence we
		// push on to the next middleware to keep the request running
		if err != nil {
			c.Next()
			return
		}

		//Create an empty interface to unmarshal our cached keys
		responseBody := map[string]interface{}{}

		//Unmarshal cached key
		json.Unmarshal(val, &responseBody)

		c.JSON(http.StatusOK, responseBody)
		// Abort other chained middlewares since we already get the response here.
		c.Abort()
	}
}

func (m RedisMiddleware) LimitPingRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Value("username").(string)
		key := fmt.Sprintf(rateLimitePrefix, username)
		_, _, allow := m.limiter.AllowMinute(key, 2)
		m.redisClient.IncrBy(key, 1)
		if !allow {
			// Handle rate limit exceeded error
			httpServer.TooManyRequests(c)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
