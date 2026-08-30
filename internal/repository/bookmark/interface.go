package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
	"gorm.io/gorm"
)

//go:generate mockery --name BookmarkRepository --filename bookmark_repository.go
type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, userId, url, description, code string) (*entity.Bookmark, error)
	GetBookmarksByUserId(ctx context.Context, userId string, limit, offset int, sort string) ([]*entity.Bookmark, int64, error)
	UpdateBookmark(ctx context.Context, userId, bookmarkId, url, description string) error
	DeleteBookmark(ctx context.Context, userId, bookmarkId string) error
}

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepository{db}
}
