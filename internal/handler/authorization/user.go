package authorization

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrorExtractClaims = errors.New("error extract token")
	ClaimsNotFound     = errors.New("Claims extract token")
)

func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenClaims, exists := c.Get("claims")

	if !exists {
		return nil, ClaimsNotFound
	}

	if _, ok := tokenClaims.(jwt.MapClaims); !ok {
		return nil, ErrorExtractClaims
	}
	return tokenClaims.(jwt.MapClaims), nil
}
