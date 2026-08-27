package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id          string    `json:"id" gorm:"type:uuid;primaryKey"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email" gorm:"type:text;not null"`
	Password    string    `json:"password" gorm:"type:text;not null"`
	UserName    string    `json:"user_name" gorm:"type:text;unique"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Bookmarks   []Bookmark
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if u.Id == "" {
		u.Id = uuid.NewString()
	}
	return
}
