package bookmark_service

import (
	"context"

	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/response"
)

var SuccessCreateNewBookmark = "success create new bookmark"

func (s *bookmarkService) NewBookmark(ctx context.Context, bookmark *bookmark_model.NewBookmarkRequest) *api.Response[bookmark_model.BookmarkInfo] {
	code := s.keyGenerator.GenerateRandomCode(10)
	bookmarkCreate, err := s.bookmarkRepository.CreateBookmark(ctx, bookmark.Url, bookmark.Description, code)
	res := &api.Response[bookmark_model.BookmarkInfo]{
		Message: SuccessCreateNewBookmark,
	}
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		return res
	}
	res.Data = domain_model.ToBookmarkInfo(bookmarkCreate)
	return res
}
