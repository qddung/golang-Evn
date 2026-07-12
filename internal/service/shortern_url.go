package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/homework/lab/internal/repository"
	"github.com/homework/lab/pkg/helpers"
	"github.com/redis/go-redis/v9"
)

type ShorternUrl interface {
	ShortenUrlShortenUrl(ctx context.Context, url string, exp int64) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

type shorternUrl struct {
	generatorRandom helpers.KeyGenerator
	repository      repository.URLStorage
}

// NewShorternUrl new shortern url
func NewShorternUrl(repository repository.URLStorage, generator helpers.KeyGenerator) ShorternUrl {
	return &shorternUrl{generator, repository}
}

// ShortenUrl shortern url
func (s *shorternUrl) ShortenUrlShortenUrl(ctx context.Context, url string, exp int64) (string, error) {
	randomCode := s.generatorRandom.GenerateRandomCode(6)
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
		return "", err
	}

	return randomCode, nil
}

var ErrCodeDoesntExist = errors.New("code does not exist")

// GetLinkFromCode return the original from shorten code
func (s *shorternUrl) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	link, err := s.repository.GetURL(ctx, code)
	if errors.Is(err, redis.Nil) {
		// try - catch none exit code exception
		return "", ErrCodeDoesntExist
	}

	return link, err
}
