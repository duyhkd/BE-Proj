package model

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
