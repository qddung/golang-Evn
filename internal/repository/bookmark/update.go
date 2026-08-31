package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// UpdateBookmark Repository
func (b *bookmarkRepository) UpdateBookmark(ctx context.Context, userId, bookmarkId, url, description string) error {
	bookmark := &entity.Bookmark{}
	err := b.db.WithContext(ctx).Model(&entity.Bookmark{}).Where("user_id = ? AND id = ?", userId, bookmarkId).First(bookmark).Error
	if err != nil {
		return err
	}

	bookmark.Url = url
	bookmark.Description = description
	return b.db.WithContext(ctx).Save(bookmark).Error
}
