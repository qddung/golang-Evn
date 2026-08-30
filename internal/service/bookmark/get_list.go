package bookmark_service

import (
	"context"

	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
)

var SuccessGetBookmarks = "success get bookmarks"

func (s *bookmarkService) GetBookmarks(ctx context.Context, userId string, query *bookmark_model.GetBookmarksQuery) (*api.PaginatedResponse[bookmark_model.BookmarkInfo], error) {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	sort := query.Sort

	bookmarks, total, err := s.bookmarkRepository.GetBookmarksByUserId(ctx, userId, limit, offset, sort)
	if err != nil {
		return nil, err
	}

	bookmarkInfos := make([]bookmark_model.BookmarkInfo, 0, len(bookmarks))
	for _, b := range bookmarks {
		bookmarkInfos = append(bookmarkInfos, *domain_model.ToBookmarkInfo(b))
	}

	res := &api.PaginatedResponse[bookmark_model.BookmarkInfo]{
		Data: bookmarkInfos,
		Pagination: api.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
	return res, nil
}

