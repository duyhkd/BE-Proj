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
	// Redis
	redisClient := redis.NewClient(&redis.Options{
		// Container name + port since we are using docker
		Addr: "localhost:6379",
	})
	// redisMiddleware := middleware.NewRedisMiddleware(redisClient.RedisClient)

	// Http Client
	app := gin.New()
	router := app.Group("/api")
	httpHandler := handlers.NewHandler(redisClient)

	router.GET("/ping", httpHandler.Ping)
	// router.GET("/", httpServer.GetRoot)
	router.POST("/login", httpHandler.Login)
	router.POST("/signup", httpHandler.SignUp)

	authorized := router.Group("")
	authorized.Use(middleware.TokenVerificationMiddleware())
	authorized.GET("/users", httpHandler.GetUserDetails)

	authorized.PATCH("/user/updatedetails", httpHandler.UpdateUserDetails)
	authorized.PATCH("/user/updatephoto", httpHandler.UpdateUserProfilePhoto)
	authorized.POST("/post", httpHandler.MakePost)
	authorized.POST("/comment", httpHandler.MakeComment)
	authorized.POST("/post/remove", httpHandler.RemovePost)
	authorized.POST("/post/update", httpHandler.EditPost)
	authorized.POST("/post/like", httpHandler.LikePost)

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
