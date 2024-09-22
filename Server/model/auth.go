package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string    `json:"username"`
	Expire   time.Time `json: "expire"`
	jwt.RegisteredClaims
}
