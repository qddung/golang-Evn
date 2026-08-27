package bookmark_handler

import (
	"github.com/gin-gonic/gin"
	bookmark_service "github.com/homework/lab/internal/service/bookmark"
)

type BookmarkHandler interface {
	CreateBookmark(c *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark_service.BookmarkService
}

func NewBookmarkHandler(svc bookmark_service.BookmarkService) BookmarkHandler {
	return &bookmarkHandler{svc}
}
