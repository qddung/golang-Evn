package bookmark_repository

import (
	"context"

	"github.com/homework/lab/internal/models/entity"
)

// CreateBookmark Repository
func (b *bookmarkRepository) CreateBookmark(ctx context.Context, url, description, code string) (*entity.Bookmark, error) {
	bk := &entity.Bookmark{Url: url, Description: description, Code: code}
	if err := b.db.WithContext(ctx).Create(bk).Error; err != nil {
		return nil, err
	}
	return bk, nil
}
