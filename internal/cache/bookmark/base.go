package bookmark_cache

import (
	"context"
	"errors"
	"fmt"

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

func (b *BookmarkCacheInstance) deleteGroupKey(ctx context.Context, userId string) error {
	err := b.cache.DeleteGroup(ctx, GetGroupKey(userId))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete group from cache in bookmarkCache.deleteGroupKey")
		return err
	}
	return nil
}
