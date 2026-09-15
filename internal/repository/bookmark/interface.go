package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
	base62_helper "github.com/homework/lab/pkg/helpers/base62"
	"gorm.io/gorm"
)

//go:generate mockery --name BookmarkRepository --filename bookmark_repository.go
type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, userId, url, description, prefix string) (*entity.Bookmark, error)
	GetBookmarksByUserId(ctx context.Context, userId string, limit, offset int, sort string) ([]*entity.Bookmark, int64, error)
	UpdateBookmark(ctx context.Context, userId, bookmarkId, url, description string) error
	DeleteBookmark(ctx context.Context, userId, bookmarkId string) error
	FindBookmarkByCode(ctx context.Context, code string) (*entity.Bookmark, error)
}

type bookmarkRepository struct {
	db            *gorm.DB
	codeGenerator base62_helper.Base62Helper
}

func NewBookmarkRepository(db *gorm.DB, codeGenerator base62_helper.Base62Helper) BookmarkRepository {
	return &bookmarkRepository{db, codeGenerator}
}
