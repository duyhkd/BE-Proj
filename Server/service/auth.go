package service

import (
	"Server/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(model.GetAppConfigs().SignedSecretKey)

func CreateToken(username string) (string, error) {
	// var existingToken model.Token
	// result := db.DB.Where("user_name = ?", username).Where("expire_at > ?", time.Now()).First(&existingToken)
	// if result.RowsAffected > 0 {
	// 	return existingToken.Token, nil
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"expire":   time.Now().Add(time.Hour * 24),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsTokenValid(tokenStr string) bool {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
