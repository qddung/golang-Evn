package infrastructure

import jwt_pkg "github.com/homework/lab/pkg/jwt"

func CreateJwtProvider() (jwt_pkg.JwtGenerator, jwt_pkg.JwtValidator) {
	jwtGenerator := jwt_pkg.NewJWTGenerator("./privatekey.pem")
	jwtValidator := jwt_pkg.NewJWTValidator("./publickey.pem")
	return jwtGenerator, jwtValidator
}
