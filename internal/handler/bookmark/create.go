package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	_ "github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
)

var SuccessCreateBookmark = "success create bookmark"

// CreateBookmark
// @Summary Create bookmark
// @Tags bookmark
// @Accept json
// @Produce json
// @Param input body bookmark_model.NewBookmarkRequest true "Input required"
// @Success 200 {object} api.Response[bookmark_model.BookmarkInfo]
// @Failure      400  {object}  api.Response[bookmark_model.BookmarkInfo]
// @Failure      500  {object}  api.Response[bookmark_model.BookmarkInfo]
// @Router /v1/bookmarks [post]
// @Security BearerAuth
func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	request, err := request_ultils.ModelBindValidation[bookmark_model.NewBookmarkRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ToDataResponse[*bookmark_model.BookmarkInfo](err))
		return
	}

	userId, err := authorization.GetSubjectFromClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}

	bookmark, err := h.svc.NewBookmark(c, userId, request)

	res := &api.Response[bookmark_model.BookmarkInfo]{}
	if err != nil {
		res.Message = response.ErrorHandling(err).Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	res.Data = bookmark
	res.Message = SuccessCreateBookmark
	c.JSON(http.StatusOK, res)
}
