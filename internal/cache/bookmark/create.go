package bookmark_cache

import (
	"context"

	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
)

// NewBookmark create new bookmark and refresh cache
func (b *BookmarkCacheInstance) NewBookmark(ctx context.Context, userId string, bookmark *bookmark_model.NewBookmarkRequest) (*bookmark_model.BookmarkInfo, error) {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return nil, err
	}
	return b.service.NewBookmark(ctx, userId, bookmark)
}
