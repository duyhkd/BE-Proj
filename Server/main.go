package main

import (
	"Server/db"
	"Server/httpServer"
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

	http.HandleFunc("/", httpServer.GetRoot)
	http.HandleFunc("/signup", httpServer.SignUp)
	http.HandleFunc("/users", httpServer.GetUserDetails)
	http.HandleFunc("/user/updatedetails", httpServer.UpdateUserDetails)
	http.HandleFunc("/user/updatephoto", httpServer.UpdateUserProfilePhoto)

	err = http.ListenAndServe(":3333", nil)
	log.Printf("starting server on port 3333")
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
