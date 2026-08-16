package user_service

import (
	"context"
	"errors"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/repository/user"
	"github.com/homework/lab/pkg/helpers"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
)

//go:generate mockery --name=UserService --filename=user_service_mock.go --outpkg=mocks
type UserService interface {
	Register(ctx context.Context, regiterInput userModel.UserRegister) (*userModel.UserInfo, error)
	Login(ctx context.Context, loginInput userModel.UserLogin) (string, error)
}

type userService struct {
	userRepository user.UserRepository
	hasher         helpers.HashHelper
	jwt            jwt_pkg.JwtGenerator
}

func NewUserService(userRepository user.UserRepository, hasher helpers.HashHelper, jwt jwt_pkg.JwtGenerator) UserService {
	return &userService{userRepository: userRepository, hasher: hasher, jwt: jwt}
}

var ServiceErr = map[string]error{
	"EmailExistError":    errors.New("Email already exists"),
	"UserNameExistError": errors.New("UserName already exists"),
	"GenerateTokenError": errors.New("Generate token error"),
}

func CheckErrorIsServiceErr(err error) bool {
	for _, svcErr := range ServiceErr {
		if svcErr == err {
			return true
		}
	}
	return false
}
