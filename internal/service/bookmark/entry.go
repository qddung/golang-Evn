package bookmark_service

import (
	"context"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	bookmark_repository "github.com/homework/lab/internal/repository/bookmark"
	"github.com/homework/lab/pkg/helpers"
)

//go:generate mockery --name BookmarkService --output mocks
type BookmarkService interface {
	NewBookmark(ctx context.Context, userId string, bookmark *bookmark_model.NewBookmarkRequest) (*bookmark_model.BookmarkInfo, error)
	GetBookmarks(ctx context.Context, userId string, query *bookmark_model.GetBookmarksQuery) (*api.PaginatedResponse[bookmark_model.BookmarkInfo], error)
	UpdateBookmark(ctx context.Context, request *bookmark_model.UpdateBookmarkRequest, userId, bookmarkId string) error
	DeleteBookmark(ctx context.Context, userId, bookmarkId string) error
}

// BookmarkService struct
type bookmarkService struct {
	bookmarkRepository bookmark_repository.BookmarkRepository
	keyGenerator       helpers.KeyGenerator
}

// Initialize BookmarkService instance
func NewBookmarkService(bookmarkRepository bookmark_repository.BookmarkRepository, keyGenerator helpers.KeyGenerator) BookmarkService {
	return &bookmarkService{bookmarkRepository, keyGenerator}
}
