package entity

import (
	"github.com/homework/lab/internal/models/base"
)

type User struct {
	base.Base
	DisplayName string `json:"display_name"`
	Email       string `json:"email" gorm:"type:text;not null"`
	Password    string `json:"password" gorm:"type:text;not null"`
	UserName    string `json:"user_name" gorm:"type:text;unique"`
	Bookmarks   []Bookmark
}
