package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/response"
)

func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	var query bookmark_model.GetBookmarksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)

	res, err := h.svc.GetBookmarks(c, userId, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
