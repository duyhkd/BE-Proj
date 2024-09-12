package httpServer

import (
	"Server/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Hello World, This is my website!\n")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		MethodNotAllowed(w)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 1 {
		BadRequest(w, "Username is missing!")
		return
	}
	if len(password) < 8 {
		BadRequest(w, "Password is too weak :D")
		return
	}

	newUser := model.User{
		UserName: username,
		Password: password,
	}

	users := model.GetUsers()
	_, ok := users[username]

	// Existing user not found
	if !ok {
		model.AddUser(newUser)
		Ok(w, fmt.Sprintf("Success fully signed up user: %s", newUser.UserName))
	} else {
		BadRequest(w, "User already signed up!")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		MethodNotAllowed(w)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, ok := model.GetUser(username)

	// Existing user not found
	if !ok || user.Password != password {
		BadRequest(w, "Credentials doesn't match or user not exist")
	} else {
		Ok(w, "Success fully logged in!")
	}
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// Need to save session/authenticated state. For now just allow anonymous
	if r.Method != "GET" {
		MethodNotAllowed(w)
	}
	username := r.URL.Query().Get("username")
	user, ok := model.GetUser(username)

	// Existing user not found
	if !ok {
		BadRequest(w, "User does not exist")
	} else {
		cleanedUser := user.AsCleanedUser()
		jsonString, _ := json.Marshal(cleanedUser)
		Ok(w, string(jsonString))
	}
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	// Need to save session/authenticated state. For now just allow anonymous
	if r.Method != "POST" {
		MethodNotAllowed(w)
	}

	username := r.URL.Query().Get("username")
	user, ok := model.GetUser(username)

	var updatedUser model.CleanedUser
	json.NewDecoder(r.Body).Decode(&updatedUser)

	// Existing user not found
	if !ok {
		BadRequest(w, "User does not exist")
	} else {
		cleanedUser := user.UpdateDetails(updatedUser)
		jsonString, _ := json.Marshal(cleanedUser)
		Ok(w, string(jsonString))
	}
}

func UpdateUserProfilePhoto(w http.ResponseWriter, r *http.Request) {
	// Need to save session/authenticated state. For now just allow anonymous
	if r.Method != "POST" {
		MethodNotAllowed(w)
	}

	err := r.ParseMultipartForm(160 << 20) // 10 MB limit
	if err != nil {
		StatusInternalServerError(w, "Error parsing form")
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		StatusInternalServerError(w, "Error retrieving file")
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
		StatusInternalServerError(w, "Error creating file")
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		StatusInternalServerError(w, "Error saving file")
		return
	}

	username := r.URL.Query().Get("username")
	user, ok := model.GetUser(username)

	// Existing user not found
	if !ok {
		BadRequest(w, "User does not exist")
	} else {
		cleanedUser := model.CleanedUser{ProfilePhoto: filePath}
		cleanedUser = user.UpdateDetails(cleanedUser)
		jsonString, _ := json.Marshal(cleanedUser)
		Ok(w, string(jsonString))
	}
}
