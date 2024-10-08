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
)

func main() {
	// Setup AppConfigs using configs.env.json
	err := model.SetupAppConfigs()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db.Init(model.GetAppConfigs().DBConnectionString)

	app := gin.New()
	router := app.Group("/api")

	router.GET("/ping", handlers.Ping)
	// router.GET("/", httpServer.GetRoot)
	router.POST("/login", handlers.Login)
	router.POST("/signup", handlers.SignUp)

	authorized := router.Group("")
	authorized.Use(middleware.TokenVerificationMiddleware())
	authorized.GET("/users", handlers.GetUserDetails)

	authorized.PATCH("/user/updatedetails", handlers.UpdateUserDetails)
	authorized.PATCH("/user/updatephoto", handlers.UpdateUserProfilePhoto)
	authorized.POST("/post", handlers.MakePost)
	authorized.POST("/comment", handlers.MakeComment)
	authorized.POST("/post/remove", handlers.RemovePost)
	authorized.POST("/post/update", handlers.EditPost)
	authorized.POST("/post/like", handlers.LikePost)

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
