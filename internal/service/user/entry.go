package user_service

import (
	"context"
	"errors"
	"reflect"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/repository/user"
	"github.com/homework/lab/pkg/helpers"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
)

// UserService interface
//
//go:generate mockery --name=UserService --filename=user_service_mock.go --outpkg=mocks
type UserService interface {
	Register(ctx context.Context, regiterInput *userModel.UserRegister) (*userModel.UserInfo, error)
	Login(ctx context.Context, loginInput userModel.UserLogin) (string, error)
	GetUserInfo(ctx context.Context, id string) (*userModel.UserInfo, error)
	UpdateUserInfo(ctx context.Context, userId string, updateInput *userModel.UpdateUserInput) error
}

type userService struct {
	userRepository user.UserRepository
	hasher         helpers.HashHelper
	jwt            jwt_pkg.JwtGenerator
}

func NewUserService(userRepository user.UserRepository, hasher helpers.HashHelper, jwt jwt_pkg.JwtGenerator) UserService {
	return &userService{userRepository: userRepository, hasher: hasher, jwt: jwt}
}

type ServiceError struct {
	EmailExistError       error
	UserNameExistError    error
	GenerateTokenError    error
	UserNameNotExistError error
	PasswordError         error
	NotFoundUserInfo      error
}

var ServiceErr = ServiceError{
	EmailExistError:       errors.New("Email already exists"),
	UserNameExistError:    errors.New("UserName already exists"),
	GenerateTokenError:    errors.New("Generate token error"),
	UserNameNotExistError: errors.New("UserName not exists"),
	PasswordError:         errors.New("Password error"),
	NotFoundUserInfo:      errors.New("User not found"),
}

func CheckErrorIsServiceErr(err error) bool {
	val := reflect.ValueOf(ServiceErr)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.Interface() == err {
			return true
		}
	}
	return false
}
