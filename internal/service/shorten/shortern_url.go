package shorten_service

import (
	"context"

	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	bookmark_repository "github.com/homework/lab/internal/repository/bookmark"
	url_repository "github.com/homework/lab/internal/repository/shorten"
	base62_helper "github.com/homework/lab/pkg/helpers/base62"
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
	repository         url_repository.URLStorage
	bookmarkRepository bookmark_repository.BookmarkRepository
	generateCode       base62_helper.Base62Helper
}

var prefix = "shorten_url_"

// NewShorternUrl new shortern url
func NewShorternUrl(repository url_repository.URLStorage, generateCode base62_helper.Base62Helper,
	bookmarkRepository bookmark_repository.BookmarkRepository) ShorternUrl {
	return &shorternUrl{repository, bookmarkRepository, generateCode}
}

// ShortenUrl shortern url
func (s *shorternUrl) ShortenUrlShortenUrl(ctx context.Context, url string, exp int64) (string, error) {
	randomCode := s.generateCode.Encode(prefix + url)
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
