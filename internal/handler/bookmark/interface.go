package bookmark_handler

import (
	"github.com/gin-gonic/gin"
	bookmark_cache "github.com/homework/lab/internal/cache/bookmark"
)

type BookmarkHandler interface {
	CreateBookmark(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark_cache.BookmarkCache
}

func NewBookmarkHandler(svc bookmark_cache.BookmarkCache) BookmarkHandler {
	return &bookmarkHandler{svc}
}
