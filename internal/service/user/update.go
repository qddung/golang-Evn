package user_service

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/dto/api/user"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/pkg/helpers"
)

var ErrorParse = errors.New("error parse")
var UpdateUserFailed = errors.New("update user failed")
var ErrorGetUser = errors.New("error get user")

func ToUpdateUser(input *user.UpdateUserInput, hasher helpers.HashHelper) (*domain_model.UpdateUser, error) {
	pass := ""
	if input.Password != "" {
		hash, err := hasher.HashPassword(input.Password)
		if err != nil {
			return nil, ErrorParse
		}
		pass = hash
	}
	return &domain_model.UpdateUser{UserName: input.UserName, Password: pass}, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, userId string, updateInput *userModel.UpdateUserInput) error {

	_, err := s.userRepository.GetUserById(ctx, userId)
	if err == gorm.ErrRecordNotFound {
		return ServiceErr.UserNameExistError
	} else if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserById in userService.UpdateUserInfo")
		return ErrorGetUser
	}
	userUpdate := &domain_model.UpdateUser{Id: userId, UserName: updateInput.UserName, Password: updateInput.Password}
	err = s.userRepository.UpdateUser(ctx, userUpdate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to UpdateUser in userService.UpdateUserInfo")
		return UpdateUserFailed
	}
	return nil
}
