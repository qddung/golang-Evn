package bookmark_cache

import (
	"context"
	"time"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/repository/cache"
	bookmark_service "github.com/homework/lab/internal/service/bookmark"
)

var ttl = time.Hour * 24 // 1 day in seconds
// Apply only read cache strategy for bookmark cache,
// since the bookmark data is not frequently updated and can be cached for a longer period of time.
// The cache will be refreshed when the bookmark data is updated or deleted.
//
//go:generate mockery --name=BookmarkCache --output=./mocks --filename=bookmark_cache_mock.go
type BookmarkCache interface {
	// Get
	GetBookmarks(ctx context.Context, userId string, query *bookmark_model.GetBookmarksQuery) (*api.PaginatedResponse[bookmark_model.BookmarkInfo], error)
	// Modify data
	NewBookmark(ctx context.Context, userId string, bookmark *bookmark_model.NewBookmarkRequest) (*bookmark_model.BookmarkInfo, error)
	UpdateBookmark(ctx context.Context, bookmark *bookmark_model.UpdateBookmarkRequest, userId string, bookmarkId string) error
	DeleteBookmark(ctx context.Context, userId string, bookmarkId string) error
}

type BookmarkCacheInstance struct {
	service bookmark_service.BookmarkService
	cache   cache.Cache
}

func NewBookmarkCacheInstance(service bookmark_service.BookmarkService, cache cache.Cache) *BookmarkCacheInstance {
	return &BookmarkCacheInstance{service: service, cache: cache}
}
