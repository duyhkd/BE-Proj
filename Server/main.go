package main

import (
	"Server/db"
	"Server/httpServer/handlers"
	"Server/middleware"
	"Server/model"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func main() {
	// Setup AppConfigs using configs.env.json
	err := model.SetupAppConfigs()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db.Init(model.GetAppConfigs().DBConnectionString)
	redisClient := redis.NewClient(&redis.Options{
		// Container name + port since we are using docker
		Addr: "localhost:6379",
	})
	redisMiddleware := middleware.NewRedisMiddleware(redisClient)

	// Http Client
	app := gin.New()
	router := app.Group("/api")
	httpHandler := handlers.NewHandler(*redisClient)

	// router.GET("/", httpServer.GetRoot)
	router.POST("/login", httpHandler.Login)
	router.POST("/user", httpHandler.SignUp)
	router.GET("/ping/top", httpHandler.TopPing)
	router.GET("/ping/count", httpHandler.PingCount)

	authorized := router.Group("")
	authorized.Use(middleware.TokenVerificationMiddleware())

	// Move out routing this is messy
	authorized.GET("/ping", redisMiddleware.LimitPingRequest(), httpHandler.Ping)
	authorized.GET("/users", httpHandler.GetUserDetails)
	authorized.PATCH("/users/details", httpHandler.UpdateUserDetails)
	authorized.PATCH("/user/updatephoto", httpHandler.UpdateUserProfilePhoto)
	authorized.POST("/posts", httpHandler.MakePost)
	authorized.POST("posts/:postId/comment", httpHandler.MakeComment)
	authorized.DELETE("/posts/:postId", httpHandler.RemovePost)
	authorized.PUT("/posts/:postId", httpHandler.EditPost)
	authorized.POST("/posts/:postId/like", httpHandler.LikePost)

	err = app.Run(":3333")
	// err = http.ListenAndServe(":3333", mux)
	log.Printf("starting server on port 3333")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
