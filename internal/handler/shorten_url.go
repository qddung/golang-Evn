package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/service"
	"github.com/homework/lab/pkg/response"
)

type ShorternUrl interface {
	ShortenUrl(c *gin.Context)
}
type shorternURL struct {
	svc service.ShorternUrl
}

func NewShortenURL(svc service.ShorternUrl) ShorternUrl {
	return &shorternURL{svc}
}

// ----- handlers -----

// ShortenLink   Generate shorten link
// @Summary      Generate shorten url based on original url that last upto 7 days
// @Description  Generate shorten url based on original url that last upto 7 days
// @Tags         link
// @Accept       application/json
// @Produce      application/json
// @Param        input body shortenInputBody true "Input required"
// @Success      200 {object} shortenResMessage
// @Router       /v1/links/shorten [post]
func (s *shorternURL) ShortenUrl(c *gin.Context) {
	request := &shortenInputBody{}
	// serialize request
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code, err := s.svc.ShortenUrlShortenUrl(c, request.Url, request.Exp)
	if err != nil {
		log.Fatal(err.Error())
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &shortenResMessage{
		Message: "Shorten URL generated successfully!",
		Code:    code,
	})

}

// -- Models
type shortenInputBody struct {
	Url string `json:"url" binding:"url,required"`
	Exp int64  `json:"exp" binding:"required,lte=604800"`
}

type shortenResMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
