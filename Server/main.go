package main

import (
	"Server/db"
	"Server/httpServer"
	"Server/httpServer/handlers"
	"Server/middleware"
	"Server/model"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Setup AppConfigs using configs.env.json
	err := model.SetupAppConfigs()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db.Init(model.GetAppConfigs().DBConnectionString)

	mux := http.NewServeMux()

	mux.HandleFunc("/", httpServer.GetRoot)
	mux.HandleFunc("/login", httpServer.Login)
	mux.HandleFunc("/signup", httpServer.SignUp)
	mux.Handle("/users", middleware.TokenVerificationMiddleware(http.HandlerFunc(httpServer.GetUserDetails)))
	mux.Handle("/user/updatedetails", middleware.TokenVerificationMiddleware(http.HandlerFunc(httpServer.UpdateUserDetails)))
	mux.Handle("/user/updatephoto", middleware.TokenVerificationMiddleware(http.HandlerFunc(httpServer.UpdateUserProfilePhoto)))
	mux.Handle("/post", middleware.TokenVerificationMiddleware(http.HandlerFunc(handlers.MakePost)))
	mux.Handle("/comment", middleware.TokenVerificationMiddleware(http.HandlerFunc(handlers.MakeComment)))
	mux.Handle("/post/remove", middleware.TokenVerificationMiddleware(http.HandlerFunc(handlers.RemovePost)))
	mux.Handle("/post/update", middleware.TokenVerificationMiddleware(http.HandlerFunc(handlers.EditPost)))

	err = http.ListenAndServe(":3333", mux)
	log.Printf("starting server on port 3333")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
