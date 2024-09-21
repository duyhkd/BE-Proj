package model

import (
	"encoding/json"
	"os"
)

const userStoragePath = "storage/user"
const usersFilePath = userStoragePath + "/users.json"

type User struct {
	UserName     string `json:"username"`
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

func AddUser(newUser User) {
	users := GetUsers()
	users[newUser.UserName] = newUser
	jsonString, _ := json.Marshal(users)
	err := os.MkdirAll(userStoragePath, os.ModePerm)
	if err != nil {
		os.Create(usersFilePath)
	}
	os.WriteFile(usersFilePath, jsonString, os.ModePerm)
}

func GetUser(username string) (User, bool) {
	users := GetUsers()
	user, ok := users[username]
	return user, ok
}
