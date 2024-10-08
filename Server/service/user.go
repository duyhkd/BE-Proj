package service

import (
	"Server/db"
	"Server/model"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AsCleanedUser(user model.User) model.CleanedUser {
	return model.CleanedUser{
		DisplayName:  user.DisplayName,
		ProfilePhoto: user.ProfilePhoto,
		Age:          user.Age,
	}
}

func UpdateDetails(username string, cleaneduser model.CleanedUser) (model.CleanedUser, error) {
	user, _ := GetUser(username)
	if cleaneduser.DisplayName != "" {
		user.DisplayName = cleaneduser.DisplayName
	}
	if cleaneduser.ProfilePhoto != "" {
		user.ProfilePhoto = cleaneduser.ProfilePhoto
	}
	if cleaneduser.Age > 0 {
		user.Age = cleaneduser.Age
	}
	result := db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_name"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "age", "profile_photo"}),
	}).Create(&user)

	return AsCleanedUser(user), result.Error
}

func GetUsers() []model.User {
	var users []model.User
	db.DB.Find(&users)
	return users
}

func AddUser(newUser model.User) error {
	var existingUser model.User
	result := db.DB.Where("user_name = ?", newUser.UserName).First(&existingUser)

	if result.Error == nil {
		// User already exists
		return result.Error
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Other db error
		return result.Error
	}

	return db.DB.Create(newUser).Error
}

func GetUser(username string) (model.User, error) {
	var user model.User
	result := db.DB.Where("user_name = ?", username).First(&user)
	return user, result.Error
}
