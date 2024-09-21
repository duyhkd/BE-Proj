package model

import (
	"Server/db"
	"encoding/json"
	"errors"
	"os"

	"gorm.io/gorm"
)

const userStoragePath = "storage/user"
const usersFilePath = userStoragePath + "/users.json"

type User struct {
	UserName     string `json:"username" gorm:"uniqueindex"`
	Password     string `json:"password"`
	DisplayName  string `json:"displayname"`
	ProfilePhoto string `json:"profilephoto"`
	Age          int    `json:"age"`
}

type CleanedUser struct {
	DisplayName  string `json:"displayname"`
	ProfilePhoto string `json:"profilephoto"`
	Age          int    `json:"age"`
}

func (user User) AsCleanedUser() CleanedUser {
	return CleanedUser{
		DisplayName:  user.DisplayName,
		ProfilePhoto: user.ProfilePhoto,
		Age:          user.Age,
	}
}

func (user User) UpdateDetails(cleaneduser CleanedUser) CleanedUser {
	

	if cleaneduser.DisplayName != "" {
		user.DisplayName = cleaneduser.DisplayName
	}
	if cleaneduser.ProfilePhoto != "" {
		user.ProfilePhoto = cleaneduser.ProfilePhoto
	}
	if cleaneduser.Age > 0 {
		user.Age = cleaneduser.Age
	}
	users := GetUsers()
	users[user.UserName] = user
	jsonString, _ := json.Marshal(users)
	os.WriteFile(usersFilePath, jsonString, os.ModePerm)
	return cleaneduser
}

func GetUsers() map[string]User {
	users := make(map[string]User)
	fileContent, err := os.ReadFile(usersFilePath)
	if err != nil {
		os.Create(usersFilePath)
	}
	json.Unmarshal(fileContent, &users)
	return users
}

func AddUser(newUser User) error {
	var existingUser User
	result := db.DB.Where("username = ?", newUser.UserName).First(&existingUser)

	if result.Error == nil {
		// User already exists
		return errors.New("user with this email already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Other db error
		return result.Error
	}

	return db.DB.Create(newUser).Error
}

func GetUser(username string) (User, bool) {
	users := GetUsers()
	user, ok := users[username]
	return user, ok
}
