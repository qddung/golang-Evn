package bookmark_repository

import (
	"context"
	"strconv"

	"github.com/homework/lab/internal/models/entity"
)

// CreateBookmark Repository
func (b *bookmarkRepository) CreateBookmark(ctx context.Context, userId, url, description, prefix string) (*entity.Bookmark, error) {
	bk := &entity.Bookmark{UserId: userId, Url: url, Description: description}
	if err := b.db.WithContext(ctx).Create(bk).Error; err != nil {
		return nil, err
	}
	code_int := strconv.Itoa(bk.CodeInt)
	code_gen := prefix + b.codeGenerator.Encode(code_int)
	bk.Code = code_gen
	b.db.Save(bk)
	return bk, nil
}
