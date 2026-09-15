package bookmark_cache

import (
	"context"
	"encoding/json"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/rs/zerolog/log"
)

// Get Bookmarks from cache
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
		} else {
			return &returnValue, nil
		}
	}

	// Return value

	res, err := b.service.GetBookmarks(ctx, userId, query)
	if err != nil {
		return nil, err
	}
	byteResult, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Marshal result in bookmarkCache.Get")
		return res, nil
	}
	err = b.cache.Write(ctx, groupKey, fields, byteResult, ttl)
	if err != nil {
		log.Error().Err(err).Msg("Failed to Write to cache in bookmarkCache.Get")
		return res, nil
	}

	return res, nil
}
