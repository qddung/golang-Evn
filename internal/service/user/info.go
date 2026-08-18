package user_service

import (
	"context"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"gorm.io/gorm"
)

func (u *userService) GetUserInfo(ctx context.Context, id string) (*userModel.UserInfo, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return nil, ServiceErr.NotFoundUserInfo
	}
	if err != nil {
		return nil, err
	}

	userInfo := &userModel.UserInfo{}
	userInfo.PopulateInfoFromUserEntity(user)
	return userInfo, nil
}
