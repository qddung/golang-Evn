package user

import (
	"context"

	domain_model "github.com/homework/lab/internal/models/domain"
	"github.com/homework/lab/internal/models/entity"
)

func (u *userRepository) UpdateUser(ctx context.Context, user *domain_model.UpdateUser) error {

	usr := &entity.User{}

	err := u.db.Where(&entity.User{
		Id: user.Id,
	}).First(usr).Error
	if err != nil {
		return err
	}
	if user.Password != "" {
		usr.Password = user.Password
	}
	if user.Password != "" {
		user.UserName = usr.UserName
	}
	return u.db.WithContext(ctx).Update("user_name", user.UserName).Error
}
