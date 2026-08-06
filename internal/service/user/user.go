package user_service

import (
	"context"
	"errors"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/repository/user"
	"github.com/homework/lab/pkg/helpers"
	"github.com/rs/zerolog/log"
)

//go:generate mockery --name=UserService --filename=user_service_mock.go --outpkg=mocks
type UserService interface {
	Register(ctx context.Context, regiterInput userModel.UserRegister) (*userModel.UserInfo, error)
}

type userService struct {
	userRepository user.UserRepository
	hasher         helpers.HashHelper
}

func NewUserService(userRepository user.UserRepository, hasher helpers.HashHelper) UserService {
	return &userService{userRepository: userRepository, hasher: hasher}
}

var EmailExistError = errors.New("Email already exists")
var UserNameExistError = errors.New("UserName already exists")

func (u *userService) Register(ctx context.Context, regiterInput userModel.UserRegister) (*userModel.UserInfo, error) {

	// check user exist
	userWithUserName, err := u.userRepository.GetUserByUserName(ctx, regiterInput.UserName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserByUserName in userService.Register")
		return nil, err
	} else if userWithUserName != nil {
		return nil, UserNameExistError
	}

	userWithEmail, err := u.userRepository.GetUserByEmail(ctx, regiterInput.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserByEmail in userService.Register")
		return nil, err
	} else if userWithEmail != nil {
		return nil, EmailExistError
	}

	hashPassword, err := u.hasher.HashPassword(regiterInput.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to HashPassword in userService.Register")
		return nil, err
	}
	// create service
	entityUser := &entity.User{
		DisplayName: regiterInput.DisplayName,
		Email:       regiterInput.Email,
		Password:    hashPassword,
		UserName:    regiterInput.UserName,
	}

	errCreateUser := u.userRepository.CreateUser(ctx, entityUser)

	if errCreateUser != nil {
		log.Error().Err(errCreateUser).Msg("Failed to create user  in userService.Register")
		return nil, errCreateUser
	}

	// get UserInfo
	userInfo := &userModel.UserInfo{}
	userInfo.PopulateInfoFromUserEntity(entityUser)
	return userInfo, nil

}
