package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
)

func (handler *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	id := c.Params.ByName("id")

	request, err := request_ultils.ModelBindValidation[bookmark_model.UpdateBookmarkRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)
	err = handler.svc.UpdateBookmark(c, request, userId, id)
	if err != nil {
		c.JSON(response.MapErrorToHttpCode[err], gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "Update bookmark successfully"})
}
