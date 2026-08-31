package route

import (
	"github.com/gin-gonic/gin"
	bookmark_handler "github.com/homework/lab/internal/handler/bookmark"
)

func BookmarkRoute(api *gin.RouterGroup, handler bookmark_handler.BookmarkHandler) {
	api.POST("/bookmarks", handler.CreateBookmark)
	api.GET("/bookmarks", handler.GetBookmarks)
	api.PUT("/bookmarks/:id", handler.UpdateBookmark)
	api.DELETE("/bookmarks/:id", handler.DeleteBookmark)
}
