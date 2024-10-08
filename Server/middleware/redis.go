package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type RedisMiddleware struct {
	redisClient redis.Client
}

func NewRedisMiddleware(
	redisClient redis.Client,
) RedisMiddleware {
	return RedisMiddleware{
		redisClient: redisClient,
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
