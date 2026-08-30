package bookmark_service

import (
	"context"

	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/response"
)

func (s *bookmarkService) UpdateBookmark(ctx context.Context, request *bookmark_model.UpdateBookmarkRequest, userId, bookmarkId string) error {
	err := s.bookmarkRepository.UpdateBookmark(ctx, userId, bookmarkId, request.Url, request.Description)
	if err != nil {
		return response.ErrorHandling(err)
	}
	return nil
}
