package shorten_service

import (
	"context"

	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/homework/lab/constant"
	bookmark_repository "github.com/homework/lab/internal/repository/bookmark"
	url_repository "github.com/homework/lab/internal/repository/shorten"
	code_generate_helper "github.com/homework/lab/pkg/helpers/code_gen"
	"github.com/redis/go-redis/v9"
)

// ShorternUrl is an interface for shorten url
//
//go:generate mockery --name=ShorternUrl --filename=shorten_url.go --outpkg=mocks
type ShorternUrl interface {
	// ShortenUrlShortenUrl method for shorten url
	ShortenUrlShortenUrl(ctx context.Context, url string, exp int64) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shorternUrl struct {
	generatorRandom    code_generate_helper.KeyGenerator
	repository         url_repository.URLStorage
	bookmarkRepository bookmark_repository.BookmarkRepository
}

// NewShorternUrl new shortern url
func NewShorternUrl(repository url_repository.URLStorage, generator code_generate_helper.KeyGenerator, bookmarkRepository bookmark_repository.BookmarkRepository) ShorternUrl {
	return &shorternUrl{generator, repository, bookmarkRepository}
}

// ShortenUrl shortern url
func (s *shorternUrl) ShortenUrlShortenUrl(ctx context.Context, url string, exp int64) (string, error) {
	randomCode := s.generatorRandom.GenerateRandomCode(constant.ShortenCodeLength)
	res, err := s.repository.GetURL(ctx, randomCode)
	// redis exeption
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	// retry get another key if key already exist
	if res != "" {
		return s.ShortenUrlShortenUrl(ctx, url, exp)
	}

	// put key into redis with value = url
	secondDuration := time.Duration(exp) * time.Second
	err = s.repository.StoreURL(ctx, randomCode, url, secondDuration)
	if err != nil {
		log.Error().Err(err).Msg("Failed to StoreURL in shorternUrl.ShortenUrlShortenUrl")
		return "", err
	}

	return randomCode, nil
}

var ErrCodeDoesntExist = errors.New("code does not exist")

// GetLinkFromCode return the original from shorten code
func (s *shorternUrl) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	if len(code) == constant.BookmarkCodeLength {
		return getLinkFromBookmarkCode(s, ctx, code)
	}
	return getLinkFromShortenCode(s, ctx, code)
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
