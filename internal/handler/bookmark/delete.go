package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/response"
)

func (handler *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Params.ByName("id")

	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)
	err = handler.svc.DeleteBookmark(c, userId, id)
	if err != nil {
		c.JSON(response.MapErrorToHttpCode[err], gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "Delete bookmark successfully"})
}
