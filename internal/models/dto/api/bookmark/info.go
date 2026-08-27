package bookmark_model

import "time"

type BookmarkInfo struct {
	Id          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
