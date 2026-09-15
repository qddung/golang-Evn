package bookmark_cache

import (
	"context"
)

// DeleteBookmark
func (b *BookmarkCacheInstance) DeleteBookmark(ctx context.Context, userId string, bookmarkId string) error {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return err
	}
	return b.service.DeleteBookmark(ctx, userId, bookmarkId)
}
