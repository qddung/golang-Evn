package bookmark_service

import (
	"context"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	bookmark_repository "github.com/homework/lab/internal/repository/bookmark"
	"github.com/homework/lab/pkg/helpers"
)

type BookmarkService interface {
	NewBookmark(ctx context.Context, bookmark *bookmark_model.NewBookmarkRequest) *api.Response[bookmark_model.BookmarkInfo]
}

type bookmarkService struct {
	bookmarkRepository bookmark_repository.BookmarkRepository
	keyGenerator       helpers.KeyGenerator
}

func NewBookmarkService(bookmarkRepository bookmark_repository.BookmarkRepository, keyGenerator helpers.KeyGenerator) BookmarkService {
	return &bookmarkService{bookmarkRepository, keyGenerator}
}
