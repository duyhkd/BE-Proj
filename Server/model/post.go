package model

import (
	"github.com/google/uuid"
)

type Post struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserName string    `json:"username"`
	Text     string    `json:"text"`
	Comments []Comment `json:"comments" gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
	Likes    []Like    `json:"likes" gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	PostId   uuid.UUID `json:"postId" gorm:"type:uuid"`
	UserName string    `json:"username"`
	Text     string    `json:"text"`
}

type Like struct {
	UserName string    `json:"username"`
	PostId   uuid.UUID `json:"postId" gorm:"type:uuid"`
}
