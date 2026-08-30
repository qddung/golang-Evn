package entity

import (
	"github.com/homework/lab/internal/models/base"
)

// Attachment represents a user-owned file/resource in the system.
// Fields follow the project's JSON and GORM tag conventions.
type Bookmark struct {
	base.Base
	Code        string `json:"code" gorm:"not null;type:text"`
	Description string `json:"description" gorm:"type:text"`
	Url         string `json:"url" gorm:"not null;type:text;index:idx_id,unique"`
	UserId      string `json:"user_id" gorm:"not null;type:uuid;index:idx_id,unique"`
	// Navigation Attribute
	User User `gorm:foreignKey:UserId`
}
