package user_service

import (
	"context"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/models/entity"
	"github.com/rs/zerolog/log"
)

func (u *userService) Register(ctx context.Context, regiterInput userModel.UserRegister) (*userModel.UserInfo, error) {

	// check user exist
	userWithUserName, err := u.userRepository.GetUserByUserName(ctx, regiterInput.UserName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserByUserName in userService.Register")
		return nil, err
	} else if userWithUserName != nil {
		return nil, ServiceErr["UserNameExistError"]
	}

	userWithEmail, err := u.userRepository.GetUserByEmail(ctx, regiterInput.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserByEmail in userService.Register")
		return nil, err
	} else if userWithEmail != nil {
		return nil, ServiceErr["EmailExistError"]
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
