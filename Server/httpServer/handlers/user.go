package handlers

import (
	"Server/httpServer"
	"Server/model"
	"Server/service"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func (handler HttpHandler) SignUp(c *gin.Context) {
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
	if len(username) < 1 {
		httpServer.BadRequest(c, "Username is missing!")
		return
	}
	if len(password) < 8 {
		httpServer.BadRequest(c, "Password is too weak :D")
		return
	}

	newUser := model.User{
		UserName: username,
		Password: password,
	}

	err := service.AddUser(newUser)

	// Existing user not found
	if err == nil {
		httpServer.Ok(c, fmt.Sprintf("Success fully signed up user: %s", newUser.UserName))
	} else {
		httpServer.BadRequest(c, "User already signed up!")
	}
}

func (handler HttpHandler) Login(c *gin.Context) {
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
	user, err := service.GetUser(username)
	cachedToken := handler.redis.Get(fmt.Sprintf("%sToken", username))
	if cachedToken.Err() != redis.Nil {
		httpServer.Ok(c, cachedToken.Val())
		return
	}

	// Existing user not found
	if err != nil || user.Password != password {
		httpServer.BadRequest(c, "Credentials doesn't match or user not exist")
	} else {
		token, _ := service.CreateToken(username)
		handler.redis.Set(fmt.Sprintf("%sToken", username), token, 180*time.Second)
		httpServer.Ok(c, token)
	}
}

func (handler HttpHandler) GetUserDetails(c *gin.Context) {
	username := c.Query("username")
	user, err := service.GetUser(username)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(c, "User does not exist")
	} else {
		cleanedUser := service.AsCleanedUser(user)
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(c, string(jsonString))
	}
}

func (handler HttpHandler) UpdateUserDetails(c *gin.Context) {
	username := c.Query("username")

	var updatedUser model.CleanedUser
	json.NewDecoder(c.Request.Body).Decode(&updatedUser)

	cleanedUser, err := service.UpdateDetails(username, updatedUser)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	} else {
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(c, string(jsonString))
	}
}

func (handler HttpHandler) UpdateUserProfilePhoto(c *gin.Context) {
	header, err := c.FormFile("photo")
	if err != nil {
		httpServer.StatusInternalServerError(c, "Error retrieving file")
		return
	}
	header.Open()

	if header.Size > 10485760 {
		httpServer.StatusInternalServerError(c, "File exceed 10MB")
		return
	}

	file, err := header.Open()
	if err != nil {
		defer file.Close()
	}

	dir := "storage/userphotos"
	filePath := filepath.Join(dir, header.Filename)
	// Ensure the directory exists
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		httpServer.StatusInternalServerError(c, "Error creating directory")
		return
	}

	// Create or open the file
	outFile, err := os.Create(filePath)
	if err != nil {
		httpServer.StatusInternalServerError(c, "Error creating file")
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		httpServer.StatusInternalServerError(c, "Error saving file")
		return
	}

	username := c.Query("username")

	cleanedUser := model.CleanedUser{ProfilePhoto: filePath}
	cleanedUser, err = service.UpdateDetails(username, cleanedUser)

	// Existing user not found
	if err != nil {
		httpServer.BadRequest(c, err.Error())
	} else {
		jsonString, _ := json.Marshal(cleanedUser)
		httpServer.Ok(c, string(jsonString))
	}
}
