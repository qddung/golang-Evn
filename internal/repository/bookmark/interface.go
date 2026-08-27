package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, url, description, code string) (*entity.Bookmark, error)
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db}
}
