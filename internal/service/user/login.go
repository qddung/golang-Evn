package user_service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/rs/zerolog/log"
)

func (u *userService) Login(ctx context.Context, loginInput userModel.UserLogin) (string, error) {
	userWithUserName, err := u.userRepository.GetUserByUserName(ctx, loginInput.UserName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to GetUserByUserName in userService.Register")
		return "", err
	}

	if userWithUserName == nil {
		return "", ServiceErr["UserNameExistError"]
	}
	claims := jwt.MapClaims{
		"sub":       userWithUserName.Id,
		"user_name": userWithUserName.UserName,
		"email":     userWithUserName.Email,
		"iat":       time.Now().String(),
	}
	token, err := u.jwt.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("Generate Token failed in jwt.GenerateToken")
		return "", ServiceErr["GenerateTokenError"]
	}
	return token, nil
}
