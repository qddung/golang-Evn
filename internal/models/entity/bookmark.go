package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Attachment represents a user-owned file/resource in the system.
// Fields follow the project's JSON and GORM tag conventions.
type Bookmark struct {
	Id          string    `json:"id" gorm:"type:uuid;primaryKey"`
	Code        string    `json:"code" gorm:"type:text"`
	Description string    `json:"description" gorm:"type:text"`
	Url         string    `json:"url" gorm:"type:text;index:idx_id,unique"`
	UserId      string    `json:"user_id" gorm:"type:uuid;index:idx_id,unique"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Navigation Attribute
	User User `gorm:foreignKey:UserId`
}

func (a *Bookmark) BeforeCreate(tx *gorm.DB) (err error) {
	if a.Id == "" {
		a.Id = uuid.NewString()
	}
	return
}
