package model

import (
	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserName  string          `json:"username"`
	Text      string          `json:"text"`
	Comments  []Comment       `json:"comments"`
	Likes     map[string]bool `json:"likes"`
	LikeCount int             `json:"total_likes" gorm:"default:0"`
}

type Comment struct {
	Id       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	PostId   uuid.UUID `json:"postId" gorm:"type:uuid"`
	UserName string    `json:"username"`
	Text     string    `json:"text"`
}
