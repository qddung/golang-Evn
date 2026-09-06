package bookmark_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/rs/zerolog/log"
)

func GetFieldsFromQuery(page, limit int) string {
	return fmt.Sprintf("%d_%d", page, limit)
}

func GetGroupKey(userId string) string {
	return fmt.Sprintf("bookmarkGroupKey_%s", userId)
}

var ErrorMarshal = errors.New("parse result failed")
var ErrorWriteCache = errors.New("write cache failed")

func (b *BookmarkCacheInstance) GetBookmarks(ctx context.Context, userId string, query *bookmark_model.GetBookmarksQuery) (*api.PaginatedResponse[bookmark_model.BookmarkInfo], error) {
	// Read Cache
	var groupKey = GetGroupKey(userId)
	var fields = GetFieldsFromQuery(query.Page, query.Limit)
	value, err := b.cache.Read(ctx, groupKey, fields)
	var returnValue api.PaginatedResponse[bookmark_model.BookmarkInfo]
	if err == nil {
		err = json.Unmarshal([]byte(value), &returnValue)
		if err != nil {
			// delete cache if unmarshal failed
			delErr := b.cache.Delete(ctx, groupKey, fields)
			if delErr != nil {
				log.Error().Err(delErr).Msg("Failed to delete cache after unmarshal failure")
			}
		}
	}
	// Return value

	// Not exit
	if err != nil {
		res, err := b.service.GetBookmarks(ctx, userId, query)
		if err != nil {
			return nil, err
		}
		byteResult, err := json.Marshal(res)
		if err != nil {
			log.Error().Err(err).Msg("Failed to Marshal result in bookmarkCache.Get")
			return nil, ErrorMarshal
		}
		err = b.cache.Write(ctx, groupKey, fields, byteResult, ttl)
		if err != nil {
			log.Error().Err(err).Msg("Failed to Write to cache in bookmarkCache.Get")
			return nil, ErrorWriteCache
		}

		return res, nil
	}

	return &returnValue, nil
}

func (b *BookmarkCacheInstance) deleteGroupKey(ctx context.Context, userId string) error {
	err := b.cache.DeleteGroup(ctx, GetGroupKey(userId))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete group from cache in bookmarkCache.deleteGroupKey")
		return err
	}
	return nil
}

func (b *BookmarkCacheInstance) NewBookmark(ctx context.Context, userId string, bookmark *bookmark_model.NewBookmarkRequest) (*bookmark_model.BookmarkInfo, error) {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return nil, err
	}
	return b.service.NewBookmark(ctx, userId, bookmark)
}

func (b *BookmarkCacheInstance) UpdateBookmark(ctx context.Context, request *bookmark_model.UpdateBookmarkRequest, userId, bookmarkId string) error {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return err
	}
	return b.service.UpdateBookmark(ctx, request, userId, bookmarkId)
}
func (b *BookmarkCacheInstance) DeleteBookmark(ctx context.Context, userId string, bookmarkId string) error {
	err := b.deleteGroupKey(ctx, userId)
	if err != nil {
		return err
	}
	return b.service.DeleteBookmark(ctx, userId, bookmarkId)
}
