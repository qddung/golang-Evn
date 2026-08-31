package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	"github.com/homework/lab/pkg/response"
)

// DeleteBookmark
// @Summary Delete bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark id"
// @Success 200 {object} api.MessageResponse
// @Failure      400  {object}  api.MessageResponse
// @Failure      500  {object}  api.MessageResponse
// @Router /v1/bookmarks/{id} [delete]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func (handler *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	id := c.Params.ByName("id")
	// authorization
	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToMessageReposnse(authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)
	// call service
	err = handler.svc.DeleteBookmark(c, userId, id)
	res := &api.MessageResponse{
		Message: "Delete bookmark successfully",
	}
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	c.JSON(http.StatusOK, res)
}
