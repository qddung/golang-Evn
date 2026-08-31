package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
)

// UpdateBookmark
// @Summary Update bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark id"
// @Param input body bookmark_model.UpdateBookmarkRequest true "Input required"
// @Success 200 {object} api.MessageResponse
// @Failure      400  {object}  api.MessageResponse
// @Failure      500  {object}  api.MessageResponse
// @Router /v1/bookmarks/{id} [put]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func (handler *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	id := c.Params.ByName("id")
	res := &api.MessageResponse{
		Message: "Update bookmark successfully",
	}
	// model validation
	request, err := request_ultils.ModelBindValidation[bookmark_model.UpdateBookmarkRequest](c)
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	// authorizaton
	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)
	// call service
	err = handler.svc.UpdateBookmark(c, request, userId, id)
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	c.JSON(http.StatusOK, res)
}
