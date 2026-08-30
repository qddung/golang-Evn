package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// DeleteBookmark Repository
func (b *bookmarkRepository) DeleteBookmark(ctx context.Context, userId, bookmarkId string) error {
	return b.db.WithContext(ctx).Where("user_id = ? AND id = ?", userId, bookmarkId).Delete(&entity.Bookmark{}).Error
}
