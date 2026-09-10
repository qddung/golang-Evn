package shorten_service

import (
	"context"
	"strings"

	bookmark_service "github.com/homework/lab/internal/service/bookmark"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ErrCodeDoesntExist = errors.New("code does not exist")

var CodeType = struct {
	Undefined int
	Shorten   int
	Bookmark  int
}{
	Undefined: 0,
	Shorten:   1,
	Bookmark:  2,
}

// GetLinkFromCode return the original from shorten code
func (s *shorternUrl) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	switch getCodeType(code) {
	case CodeType.Shorten:
		return getLinkFromShortenCode(s, ctx, code)
	case CodeType.Bookmark:
		return getLinkFromBookmarkCode(s, ctx, code)
	}
	return "", ErrCodeDoesntExist
}

func getCodeType(code string) int {
	if strings.Contains(code, bookmark_service.PrefixCreate) {
		return CodeType.Bookmark
	}

	if strings.Contains(code, prefix) {
		return CodeType.Shorten
	}
	return CodeType.Undefined
}

func getLinkFromShortenCode(s *shorternUrl, ctx context.Context, code string) (string, error) {
	link, err := s.repository.GetURL(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeDoesntExist
	}

	return link, err
}

func getLinkFromBookmarkCode(s *shorternUrl, ctx context.Context, code string) (string, error) {
	bookmark, err := s.bookmarkRepository.FindBookmarkByCode(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeDoesntExist
	}
	return bookmark.Url, err
}
