package handlers

import (
	"Server/httpServer"
	"Server/model"
	"Server/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 1 {
		httpServer.BadRequest(w, "Username is missing!")
		return
	}
	if len(password) < 8 {
		httpServer.BadRequest(w, "Password is too weak :D")
		return
	}

	newUser := model.User{
		UserName: username,
		Password: password,
	}

	err := service.AddUser(newUser)

	// Existing user not found
	if err == nil {
		httpServer.Ok(w, fmt.Sprintf("Success fully signed up user: %s", newUser.UserName))
	} else {
		httpServer.BadRequest(w, "User already signed up!")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := service.GetUser(username)

	// Existing user not found
	if err != nil || user.Password != password {
		httpServer.BadRequest(w, "Credentials doesn't match or user not exist")
	} else {
		token, _ := service.CreateToken(username)
		httpServer.Ok(w, token)
	}
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpServer.MethodNotAllowed(w)
	}
	username := r.URL.Query().Get("username")
	user, err := service.GetUser(username)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(w, "User does not exist")
	} else {
		cleanedUser := service.AsCleanedUser(user)
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(w, string(jsonString))
	}
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	username := r.URL.Query().Get("username")

	var updatedUser model.CleanedUser
	json.NewDecoder(r.Body).Decode(&updatedUser)

	cleanedUser, err := service.UpdateDetails(username, updatedUser)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	} else {
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(w, string(jsonString))
	}
}

func UpdateUserProfilePhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpServer.MethodNotAllowed(w)
	}

	err := r.ParseMultipartForm(160 << 20) // 10 MB limit
	if err != nil {
		httpServer.StatusInternalServerError(w, "Error parsing form")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		httpServer.StatusInternalServerError(w, "Error retrieving file")
		return
	}
	defer file.Close()

	dir := "storage/userphotos"
	filePath := filepath.Join(dir, header.Filename)
	// Ensure the directory exists
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Create or open the file
	outFile, err := os.Create(filePath)
	if err != nil {
		httpServer.StatusInternalServerError(w, "Error creating file")
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		httpServer.StatusInternalServerError(w, "Error saving file")
		return
	}

	username := r.URL.Query().Get("username")

	cleanedUser := model.CleanedUser{ProfilePhoto: filePath}
	cleanedUser, err = service.UpdateDetails(username, cleanedUser)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(w, err.Error())
	} else {
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(w, string(jsonString))
	}
}
