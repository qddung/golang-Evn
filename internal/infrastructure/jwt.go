package infrastructure

import (
	"github.com/homework/lab/constant"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
)

func CreateJwtProvider() (jwt_pkg.JwtGenerator, jwt_pkg.JwtValidator) {
	jwtGenerator := jwt_pkg.NewJWTGenerator(constant.PrivateKeyPath)
	jwtValidator := jwt_pkg.NewJWTValidator(constant.PublicKeyPath)
	return jwtGenerator, jwtValidator
}
