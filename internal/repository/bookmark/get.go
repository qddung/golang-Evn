package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// GetBookmarksByUserId Repository return list bookmark, total bookmark, error
func (b *bookmarkRepository) GetBookmarksByUserId(ctx context.Context, userId string, limit, offset int, sort string) ([]*entity.Bookmark, int64, error) {
	var bookmarks []*entity.Bookmark
	var total int64

	query := b.db.WithContext(ctx).Model(&entity.Bookmark{}).Where("user_id = ?", userId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.Bookmark{}, 0, nil
	}

	if sort != "" {
		query = query.Order(sort)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

func (b *bookmarkRepository) FindBookmarkByCode(ctx context.Context, code string) (*entity.Bookmark, error) {
	var bookmark entity.Bookmark
	if err := b.db.WithContext(ctx).Where("code = ?", code).First(&bookmark).Error; err != nil {
		return nil, err
	}
	return &bookmark, nil
}
